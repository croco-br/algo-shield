package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/tokenrevoke"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user operations needed by auth
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, email, name, passwordHash string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error
}

type Service struct {
	userService        UserService
	jwtSecret          string
	jwtExpiry          time.Duration
	tokenRevokeService *tokenrevoke.Service
}

// NewService creates a new auth service with dependency injection
// Follows Dependency Inversion Principle - receives interface, not concrete type
func NewService(cfg *config.Config, userService UserService, tokenRevokeService *tokenrevoke.Service) *Service {
	return &Service{
		userService:        userService,
		jwtSecret:          cfg.Auth.JWTSecret,
		jwtExpiry:          time.Duration(cfg.Auth.JWTExpirationHours) * time.Hour,
		tokenRevokeService: tokenRevokeService,
	}
}

// RegisterUser handles user registration with password hashing
func (s *Service) RegisterUser(ctx context.Context, email, name, password string) (*models.User, string, error) {
	// Check if user already exists
	existingUser, err := s.userService.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, "", fmt.Errorf("user with this email already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user via user service
	user, err := s.userService.CreateUser(ctx, email, name, string(passwordHash))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// LoginUser handles user login with password verification
func (s *Service) LoginUser(ctx context.Context, email, password string) (*models.User, string, error) {
	// Get user by email with password
	user, err := s.userService.GetUserByEmailWithPassword(ctx, email)
	if err != nil {
		// Use safe error message - don't reveal if email exists
		return nil, "", apierrors.InvalidCredentials()
	}

	// Verify password
	if user.PasswordHash == nil {
		return nil, "", apierrors.InvalidCredentials()
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", apierrors.InvalidCredentials()
	}

	if !user.Active {
		return nil, "", apierrors.UserInactive()
	}

	// Update last login (non-critical, don't fail login if this fails)
	now := time.Now()
	if err := s.userService.UpdateLastLogin(ctx, user.ID, &now); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", apierrors.InternalError("Failed to generate authentication token")
	}

	return user, token, nil
}

// GenerateJWT generates a JWT token for a user with proper claims
// Following JWT best practices: iat (issued at), exp (expiration), user_id, email, name
func (s *Service) GenerateJWT(user *models.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"name":    user.Name,
		"iat":     now.Unix(),       // Issued at
		"exp":     expiresAt.Unix(), // Expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user
// Performs comprehensive validation: signature, expiration, revocation, user status
func (s *Service) ValidateToken(tokenString string) (*models.User, error) {
	// Parse and validate JWT signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierrors.TokenInvalid()
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		// Check if token is expired
		if err.Error() == "token is expired" || err.Error() == "Token is expired" {
			return nil, apierrors.TokenExpired()
		}
		return nil, apierrors.TokenInvalid()
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apierrors.TokenInvalid()
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, apierrors.TokenInvalid()
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apierrors.TokenInvalid()
	}

	// Check if token is revoked (individual token blacklist)
	ctx, cancel := context.WithTimeout(context.Background(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	isRevoked, err := s.tokenRevokeService.IsTokenRevoked(ctx, tokenString)
	if err != nil {
		// Log error but don't fail validation (graceful degradation)
		log.Printf("Failed to check token revocation status: %v", err)
	} else if isRevoked {
		return nil, apierrors.TokenRevoked()
	}

	// Check if all user tokens are revoked (e.g., password change)
	userTokensRevoked, err := s.tokenRevokeService.IsUserTokensRevoked(ctx, userID.String())
	if err != nil {
		// Log error but don't fail validation (graceful degradation)
		log.Printf("Failed to check user tokens revocation status: %v", err)
	} else if userTokensRevoked {
		return nil, apierrors.TokenRevoked()
	}

	// Get user from database
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apierrors.NotFound("User")
	}

	// Check if user is active
	if !user.Active {
		return nil, apierrors.UserInactive()
	}

	return user, nil
}

// LogoutUser revokes the current token
func (s *Service) LogoutUser(ctx context.Context, tokenString string) error {
	// Parse token to get expiration time
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		// Even if token is invalid/expired, we don't return error on logout
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// Get expiration time from claims
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil
	}

	expiresAt := time.Unix(int64(exp), 0)

	// Revoke the token
	return s.tokenRevokeService.RevokeToken(ctx, tokenString, expiresAt)
}

// RevokeAllUserTokens revokes all tokens for a user (e.g., on password change or account deactivation)
func (s *Service) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	return s.tokenRevokeService.RevokeAllUserTokens(ctx, userID.String(), s.jwtExpiry)
}
