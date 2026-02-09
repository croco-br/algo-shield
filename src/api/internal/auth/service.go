package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	RegisterUser(ctx context.Context, email, name, password string) (*models.User, string, string, error)
	LoginUser(ctx context.Context, email, password string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error)
	ValidateToken(ctx context.Context, tokenString string) (*models.User, error)
	LogoutUser(ctx context.Context, tokenString string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	GenerateJWT(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
}

// UserService defines the interface for user operations needed by auth
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, email, name, passwordHash string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error
}

// TokenRevokeService defines the interface for token revocation operations
type TokenRevokeService interface {
	RevokeToken(ctx context.Context, token string, expiresAt time.Time) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
	RevokeAllUserTokens(ctx context.Context, userID string, tokenExpiry time.Duration) error
	IsUserTokensRevoked(ctx context.Context, userID string) (bool, error)
	IsTokenOrUserRevoked(ctx context.Context, token string, userID string) (bool, error)
}

type Service struct {
	userService        UserService
	jwtSecret          string
	jwtExpiry          time.Duration
	jwtRefreshExpiry   time.Duration
	tokenRevokeService TokenRevokeService
}

// NewService creates a new auth service with dependency injection
// Follows Dependency Inversion Principle - receives interface, not concrete type
func NewService(cfg *config.Config, userService UserService, tokenRevokeService TokenRevokeService) *Service {
	return &Service{
		userService:        userService,
		jwtSecret:          cfg.Auth.JWTSecret,
		jwtExpiry:          time.Duration(cfg.Auth.JWTExpirationHours) * time.Hour,
		jwtRefreshExpiry:   time.Duration(cfg.Auth.JWTRefreshExpirationHours) * time.Hour,
		tokenRevokeService: tokenRevokeService,
	}
}

// RegisterUser handles user registration with password hashing
func (s *Service) RegisterUser(ctx context.Context, email, name, password string) (*models.User, string, string, error) {
	// Check if user already exists
	existingUser, err := s.userService.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, "", "", fmt.Errorf("user with this email already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user via user service
	user, err := s.userService.CreateUser(ctx, email, name, string(passwordHash))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT access token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user, token, refreshToken, nil
}

// LoginUser handles user login with password verification
func (s *Service) LoginUser(ctx context.Context, email, password string) (*models.User, string, string, error) {
	// Get user by email with password
	user, err := s.userService.GetUserByEmailWithPassword(ctx, email)
	if err != nil {
		// Use safe error message - don't reveal if email exists
		return nil, "", "", apierrors.InvalidCredentials()
	}

	// Verify password
	if user.PasswordHash == nil {
		return nil, "", "", apierrors.InvalidCredentials()
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", "", apierrors.InvalidCredentials()
	}

	if !user.Active {
		return nil, "", "", apierrors.UserInactive()
	}

	// Update last login (non-critical, don't fail login if this fails)
	now := time.Now()
	if err := s.userService.UpdateLastLogin(ctx, user.ID, &now); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate JWT access token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", "", apierrors.InternalError("Failed to generate authentication token")
	}

	// Generate refresh token
	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", apierrors.InternalError("Failed to generate refresh token")
	}

	return user, token, refreshToken, nil
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
		"type":    "access",         // Token type
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token with longer expiration
// Refresh tokens are used to obtain new access tokens without re-authentication
func (s *Service) GenerateRefreshToken(user *models.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtRefreshExpiry)

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"iat":     now.Unix(),       // Issued at
		"exp":     expiresAt.Unix(), // Expiration time
		"type":    "refresh",        // Token type
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// RefreshToken validates a refresh token and generates new access + refresh tokens (rotation)
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error) {
	// Parse and validate JWT signature
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierrors.TokenInvalid()
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, "", "", apierrors.TokenExpired()
		}
		return nil, "", "", apierrors.TokenInvalid()
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, "", "", apierrors.TokenInvalid()
	}

	// Verify this is a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, "", "", apierrors.NewAPIError(apierrors.ErrUnauthorized, "Invalid token type")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, "", "", apierrors.TokenInvalid()
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, "", "", apierrors.TokenInvalid()
	}

	// Check if token or user tokens are revoked (single Redis round-trip)
	isRevoked, err := s.tokenRevokeService.IsTokenOrUserRevoked(ctx, refreshToken, userID.String())
	if err != nil {
		log.Printf("Failed to check refresh token revocation status: %v", err)
		return nil, "", "", apierrors.InternalError("Unable to verify token status")
	}
	if isRevoked {
		return nil, "", "", apierrors.TokenRevoked()
	}

	// Get user from database
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, "", "", apierrors.NotFound("User")
	}

	// Check if user is active
	if !user.Active {
		return nil, "", "", apierrors.UserInactive()
	}

	// Revoke the old refresh token (rotation: one-time use)
	exp, _ := claims["exp"].(float64)
	oldExpiresAt := time.Unix(int64(exp), 0)
	if err := s.tokenRevokeService.RevokeToken(ctx, refreshToken, oldExpiresAt); err != nil {
		log.Printf("Failed to revoke old refresh token: %v", err)
	}

	// Generate new access token
	newAccessToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", "", apierrors.InternalError("Failed to generate access token")
	}

	// Generate new refresh token (rotation)
	newRefreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", apierrors.InternalError("Failed to generate refresh token")
	}

	return user, newAccessToken, newRefreshToken, nil
}

// ValidateToken validates a JWT token and returns the user
// Performs comprehensive validation: signature, expiration, type, revocation, user status
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Parse and validate JWT signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apierrors.TokenInvalid()
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apierrors.TokenExpired()
		}
		return nil, apierrors.TokenInvalid()
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apierrors.TokenInvalid()
	}

	// Verify this is an access token (prevent refresh tokens from being used as Bearer tokens)
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
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

	// Check if token or user tokens are revoked (single Redis round-trip, fail-closed)
	isRevoked, err := s.tokenRevokeService.IsTokenOrUserRevoked(ctx, tokenString, userID.String())
	if err != nil {
		log.Printf("Failed to check token revocation status: %v", err)
		return nil, apierrors.InternalError("Unable to verify token status")
	}
	if isRevoked {
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
