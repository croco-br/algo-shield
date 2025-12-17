package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	userService UserService
	jwtSecret   string
	jwtExpiry   time.Duration
}

func NewService(db *pgxpool.Pool, cfg *config.Config, userService UserService) *Service {
	return &Service{
		userService: userService,
		jwtSecret:   cfg.Auth.JWTSecret,
		jwtExpiry:   time.Duration(cfg.Auth.JWTExpirationHours) * time.Hour,
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
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Verify password
	if user.PasswordHash == nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password")
	}

	if !user.Active {
		return nil, "", fmt.Errorf("user account is inactive")
	}

	// Update last login (non-critical, don't fail login if this fails)
	now := time.Now()
	if err := s.userService.UpdateLastLogin(ctx, user.ID, &now); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GenerateJWT generates a JWT token for a user
func (s *Service) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"name":    user.Name,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the user
func (s *Service) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		ctx, cancel := context.WithTimeout(context.Background(), internal.DEFAULT_TIMEOUT)
		defer cancel()
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		// Get user from database via user service
		user, err := s.userService.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if !user.Active {
			return nil, fmt.Errorf("user is inactive")
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}
