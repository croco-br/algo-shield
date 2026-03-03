package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

// Test_Service_LoginUser_WhenValidCredentials_ThenReturnsUserAndToken tests successful login
func Test_Service_LoginUser_WhenValidCredentials_ThenReturnsUserAndToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "test@example.com"
	password := "Password123!"

	// Hash password for mock user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	passwordHashStr := string(passwordHash)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:           userID,
		Email:        email,
		Name:         "Test User",
		Active:       true,
		PasswordHash: &passwordHashStr,
	}

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(expectedUser, nil)

	mockUserService.EXPECT().
		UpdateLastLogin(ctx, userID, gomock.Any()).
		Return(nil)

	// Act
	user, token, _, err := service.LoginUser(ctx, email, password)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, userID, user.ID)
}

// Test_Service_LoginUser_WhenInvalidEmail_ThenReturnsInvalidCredentialsError tests safe error message
func Test_Service_LoginUser_WhenInvalidEmail_ThenReturnsInvalidCredentialsError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "nonexistent@example.com"
	password := "Password123!"

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(nil, errors.New("user not found"))

	// Act
	user, token, _, err := service.LoginUser(ctx, email, password)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Verify it returns the safe error message (doesn't reveal if email exists)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrInvalidCredentials, apiErr.Code)
	assert.Equal(t, "Invalid email or password", apiErr.Message)
}

// Test_Service_LoginUser_WhenWrongPassword_ThenReturnsInvalidCredentialsError tests safe error for wrong password
func Test_Service_LoginUser_WhenWrongPassword_ThenReturnsInvalidCredentialsError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "test@example.com"
	correctPassword := "Password123!"
	wrongPassword := "wrongpassword"

	// Hash correct password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	require.NoError(t, err)
	passwordHashStr := string(passwordHash)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:           userID,
		Email:        email,
		Name:         "Test User",
		Active:       true,
		PasswordHash: &passwordHashStr,
	}

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(expectedUser, nil)

	// Act
	user, token, _, err := service.LoginUser(ctx, email, wrongPassword)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	// Verify it returns the safe error message (same as non-existent email)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrInvalidCredentials, apiErr.Code)
	assert.Equal(t, "Invalid email or password", apiErr.Message)
}

// Test_Service_LoginUser_WhenUserInactive_ThenReturnsUserInactiveError tests inactive user handling
func Test_Service_LoginUser_WhenUserInactive_ThenReturnsUserInactiveError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "test@example.com"
	password := "Password123!"

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	passwordHashStr := string(passwordHash)

	userID := uuid.New()
	inactiveUser := &models.User{
		ID:           userID,
		Email:        email,
		Name:         "Test User",
		Active:       false, // User is inactive
		PasswordHash: &passwordHashStr,
	}

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(inactiveUser, nil)

	// Act
	user, token, _, err := service.LoginUser(ctx, email, password)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrUserInactive, apiErr.Code)
}

// Test_Service_GenerateJWT_WhenValidUser_ThenReturnsTokenWithClaims tests JWT generation
func Test_Service_GenerateJWT_WhenValidUser_ThenReturnsTokenWithClaims(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Act
	tokenString, err := service.GenerateJWT(user)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse and verify token claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// Verify required claims exist
	assert.Equal(t, userID.String(), claims["user_id"])
	assert.Equal(t, user.Email, claims["email"])
	assert.Equal(t, user.Name, claims["name"])
	assert.NotNil(t, claims["iat"], "Token should have issued at (iat) claim")
	assert.NotNil(t, claims["exp"], "Token should have expiration (exp) claim")

	// Verify expiration is approximately 24 hours from now
	exp := int64(claims["exp"].(float64))
	expectedExp := time.Now().Add(24 * time.Hour).Unix()
	assert.InDelta(t, expectedExp, exp, 5, "Expiration should be ~24 hours from now")
}

// Test_Service_ValidateToken_WhenValidToken_ThenReturnsUser tests token validation
func Test_Service_ValidateToken_WhenValidToken_ThenReturnsUser(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate a valid token
	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	// Set up mock expectations before calling ValidateToken
	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(false, nil).
		Times(1)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(user, nil).
		Times(1)

	// Act
	validatedUser, err := service.ValidateToken(ctx, tokenString)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, validatedUser)
	assert.Equal(t, userID, validatedUser.ID)
	assert.Equal(t, user.Email, validatedUser.Email)
}

// Test_Service_ValidateToken_WhenInvalidToken_ThenReturnsTokenInvalidError tests invalid token
func Test_Service_ValidateToken_WhenInvalidToken_ThenReturnsTokenInvalidError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	invalidToken := "invalid.token.here"
	ctx := context.Background()

	// Act
	user, err := service.ValidateToken(ctx, invalidToken)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrTokenInvalid, apiErr.Code)
}

// Test_Service_ValidateToken_WhenRevokedToken_ThenReturnsTokenRevokedError tests revoked token
func Test_Service_ValidateToken_WhenRevokedToken_ThenReturnsTokenRevokedError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate a valid token
	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	// Mock token as revoked (combined check)
	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(true, nil)

	// Act
	validatedUser, err := service.ValidateToken(ctx, tokenString)

	// Assert
	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrTokenRevoked, apiErr.Code)
}

// Test_Service_ValidateToken_WhenUserTokensRevoked_ThenReturnsTokenRevokedError tests user-level revocation
func Test_Service_ValidateToken_WhenUserTokensRevoked_ThenReturnsTokenRevokedError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate a valid token
	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	// All user tokens are revoked (e.g., password change) — combined check returns true
	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(true, nil)

	// Act
	validatedUser, err := service.ValidateToken(ctx, tokenString)

	// Assert
	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrTokenRevoked, apiErr.Code)
}

// Test_Service_LogoutUser_WhenValidToken_ThenRevokesToken tests logout functionality
func Test_Service_LogoutUser_WhenValidToken_ThenRevokesToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate a valid token
	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		RevokeToken(ctx, tokenString, gomock.Any()).
		Return(nil)

	// Act
	err = service.LogoutUser(ctx, tokenString)

	// Assert
	require.NoError(t, err)
}

// Test_Service_LogoutUser_WhenInvalidToken_ThenDoesNotFail tests graceful logout with invalid token
func Test_Service_LogoutUser_WhenInvalidToken_ThenDoesNotFail(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	invalidToken := "invalid.token.here"

	// Act
	err := service.LogoutUser(ctx, invalidToken)

	// Assert
	// Logout should not fail even with invalid token
	require.NoError(t, err)
}

// Test_Service_RegisterUser_WhenValidData_ThenCreatesUserAndToken tests registration
func Test_Service_RegisterUser_WhenValidData_ThenCreatesUserAndToken(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	password := "Password123!"

	userID := uuid.New()
	createdUser := &models.User{
		ID:     userID,
		Email:  email,
		Name:   name,
		Active: true,
	}

	// User doesn't exist yet
	mockUserService.EXPECT().
		GetUserByEmail(ctx, email).
		Return(nil, errors.New("not found"))

	// User is created
	mockUserService.EXPECT().
		CreateUser(ctx, email, name, gomock.Any()).
		Return(createdUser, nil)

	// Act
	user, token, _, err := service.RegisterUser(ctx, email, name, password)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
}

// Test_Service_RegisterUser_WhenEmailExists_ThenReturnsError tests duplicate email handling
func Test_Service_RegisterUser_WhenEmailExists_ThenReturnsError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "existing@example.com"
	name := "New User"
	password := "Password123!"

	existingUser := &models.User{
		ID:     uuid.New(),
		Email:  email,
		Name:   "Existing User",
		Active: true,
	}

	// User already exists
	mockUserService.EXPECT().
		GetUserByEmail(ctx, email).
		Return(existingUser, nil)

	// Act
	user, token, _, err := service.RegisterUser(ctx, email, name, password)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "already exists")
}

// Test_Service_RegisterUser_WhenCreateUserFails_ThenReturnsError tests user creation failure
func Test_Service_RegisterUser_WhenCreateUserFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	password := "Password123!"

	mockUserService.EXPECT().
		GetUserByEmail(ctx, email).
		Return(nil, errors.New("not found"))

	mockUserService.EXPECT().
		CreateUser(ctx, email, name, gomock.Any()).
		Return(nil, errors.New("database error"))

	user, token, _, err := service.RegisterUser(ctx, email, name, password)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "failed to create user")
}

// Test_Service_LoginUser_WhenPasswordHashIsNil_ThenReturnsInvalidCredentials tests nil password hash
func Test_Service_LoginUser_WhenPasswordHashIsNil_ThenReturnsInvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	email := "test@example.com"
	password := "Password123!"

	userWithoutPassword := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         "Test User",
		Active:       true,
		PasswordHash: nil,
	}

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(userWithoutPassword, nil)

	user, token, _, err := service.LoginUser(ctx, email, password)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrInvalidCredentials, apiErr.Code)
}

// Test_Service_ValidateToken_WhenInvalidSigningMethod_ThenReturnsError tests invalid signing method
func Test_Service_ValidateToken_WhenInvalidSigningMethod_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	invalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.invalid"
	ctx := context.Background()

	user, err := service.ValidateToken(ctx, invalidToken)

	require.Error(t, err)
	assert.Nil(t, user)
}

// Test_Service_ValidateToken_WhenUserNotFound_ThenReturnsError tests user not found
func Test_Service_ValidateToken_WhenUserNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	validatedUser, err := service.ValidateToken(ctx, tokenString)

	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrNotFound, apiErr.Code)
}

// Test_Service_ValidateToken_WhenUserInactive_ThenReturnsError tests inactive user during validation
func Test_Service_ValidateToken_WhenUserInactive_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	inactiveUser := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: false,
	}

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(inactiveUser, nil)

	validatedUser, err := service.ValidateToken(ctx, tokenString)

	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrUserInactive, apiErr.Code)
}

// Test_Service_RevokeAllUserTokens_WhenCalled_ThenCallsTokenRevokeService tests RevokeAllUserTokens
func Test_Service_RevokeAllUserTokens_WhenCalled_ThenCallsTokenRevokeService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	userID := uuid.New()

	mockTokenRevoke.EXPECT().
		RevokeAllUserTokens(ctx, userID.String(), 24*time.Hour).
		Return(nil)

	err := service.RevokeAllUserTokens(ctx, userID)

	require.NoError(t, err)
}

// Test_Service_RevokeAllUserTokens_WhenServiceFails_ThenReturnsError tests RevokeAllUserTokens error
func Test_Service_RevokeAllUserTokens_WhenServiceFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()
	userID := uuid.New()

	mockTokenRevoke.EXPECT().
		RevokeAllUserTokens(ctx, userID.String(), 24*time.Hour).
		Return(errors.New("redis error"))

	err := service.RevokeAllUserTokens(ctx, userID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis error")
}

// Test_Service_RefreshToken_WhenValidToken_ThenReturnsNewTokens tests successful token refresh
func Test_Service_RefreshToken_WhenValidToken_ThenReturnsNewTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), refreshToken, userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(user, nil)

	// Old token gets revoked (rotation)
	mockTokenRevoke.EXPECT().
		RevokeToken(gomock.Any(), refreshToken, gomock.Any()).
		Return(nil)

	returnedUser, newAccessToken, newRefreshToken, err := service.RefreshToken(ctx, refreshToken)

	require.NoError(t, err)
	assert.NotNil(t, returnedUser)
	assert.Equal(t, userID, returnedUser.ID)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
}

// Test_Service_RefreshToken_WhenExpiredToken_ThenReturnsTokenExpired tests expired refresh token
func Test_Service_RefreshToken_WhenExpiredToken_ThenReturnsTokenExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	// Create an expired refresh token
	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "refresh",
		"iat":     time.Now().Add(-200 * time.Hour).Unix(),
		"exp":     time.Now().Add(-1 * time.Hour).Unix(), // Already expired
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString([]byte(jwtSecret))
	require.NoError(t, err)

	ctx := context.Background()

	_, _, _, err = service.RefreshToken(ctx, expiredToken)

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrTokenExpired, apiErr.Code)
}

// Test_Service_RefreshToken_WhenAccessTokenUsedAsRefresh_ThenReturnsError tests token type enforcement
func Test_Service_RefreshToken_WhenAccessTokenUsedAsRefresh_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	user := &models.User{
		ID:     uuid.New(),
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate an access token and try to use it as a refresh token
	accessToken, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	_, _, _, err = service.RefreshToken(ctx, accessToken)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid token type")
}

// Test_Service_RefreshToken_WhenTokenRevoked_ThenReturnsTokenRevoked tests revoked token
func Test_Service_RefreshToken_WhenTokenRevoked_ThenReturnsTokenRevoked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), refreshToken, userID.String()).
		Return(true, nil)

	_, _, _, err = service.RefreshToken(ctx, refreshToken)

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrTokenRevoked, apiErr.Code)
}

// Test_Service_RefreshToken_WhenRedisError_ThenFailsClosed tests fail-closed on Redis error
func Test_Service_RefreshToken_WhenRedisError_ThenFailsClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), refreshToken, userID.String()).
		Return(false, errors.New("redis connection error"))

	_, _, _, err = service.RefreshToken(ctx, refreshToken)

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrInternalError, apiErr.Code)
}

// Test_Service_RefreshToken_WhenUserInactive_ThenReturnsUserInactive tests inactive user during refresh
func Test_Service_RefreshToken_WhenUserInactive_ThenReturnsUserInactive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	inactiveUser := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: false,
	}

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), refreshToken, userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(inactiveUser, nil)

	_, _, _, err = service.RefreshToken(ctx, refreshToken)

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrUserInactive, apiErr.Code)
}

// Test_Service_RefreshToken_WhenUserNotFound_ThenReturnsNotFound tests user not found during refresh
func Test_Service_RefreshToken_WhenUserNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), refreshToken, userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	_, _, _, err = service.RefreshToken(ctx, refreshToken)

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrNotFound, apiErr.Code)
}

// Test_Service_RefreshToken_WhenInvalidToken_ThenReturnsTokenInvalid tests invalid token
func Test_Service_RefreshToken_WhenInvalidToken_ThenReturnsTokenInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()

	_, _, _, err := service.RefreshToken(ctx, "invalid.token.here")

	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok)
	assert.Equal(t, apierrors.ErrTokenInvalid, apiErr.Code)
}

// Test_Service_RevokeRefreshToken_WhenEmptyToken_ThenReturnsNil tests empty token
func Test_Service_RevokeRefreshToken_WhenEmptyToken_ThenReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()

	err := service.RevokeRefreshToken(ctx, "")

	require.NoError(t, err)
}

// Test_Service_RevokeRefreshToken_WhenValidToken_ThenRevokesToken tests successful revocation
func Test_Service_RevokeRefreshToken_WhenValidToken_ThenRevokesToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	user := &models.User{
		ID:     uuid.New(),
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	mockTokenRevoke.EXPECT().
		RevokeToken(ctx, refreshToken, gomock.Any()).
		Return(nil)

	err = service.RevokeRefreshToken(ctx, refreshToken)

	require.NoError(t, err)
}

// Test_Service_RevokeRefreshToken_WhenInvalidToken_ThenReturnsNil tests invalid token is silently ignored
func Test_Service_RevokeRefreshToken_WhenInvalidToken_ThenReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key",
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	ctx := context.Background()

	err := service.RevokeRefreshToken(ctx, "invalid.token.here")

	require.NoError(t, err)
}

// Test_Service_ValidateToken_WhenRefreshTokenUsedAsAccess_ThenReturnsTokenInvalid tests token type enforcement
func Test_Service_ValidateToken_WhenRefreshTokenUsedAsAccess_ThenReturnsTokenInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:                 "test-secret-key",
			JWTExpirationHours:        24,
			JWTRefreshExpirationHours: 168,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	// Generate a refresh token
	refreshToken, err := service.GenerateRefreshToken(user)
	require.NoError(t, err)

	ctx := context.Background()

	// Act: try to use refresh token as access token
	validatedUser, err := service.ValidateToken(ctx, refreshToken)

	// Assert: should be rejected
	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrTokenInvalid, apiErr.Code)
}

// Test_Service_ValidateToken_WhenRedisError_ThenFailsClosed tests fail-closed on Redis error
func Test_Service_ValidateToken_WhenRedisError_ThenFailsClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := NewMockUserService(ctrl)
	mockTokenRevoke := NewMockTokenRevokeService(ctrl)

	jwtSecret := "test-secret-key"
	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: 24,
		},
	}

	service := NewService(cfg, mockUserService, mockTokenRevoke)

	userID := uuid.New()
	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	tokenString, err := service.GenerateJWT(user)
	require.NoError(t, err)

	ctx := context.Background()

	// Redis returns an error
	mockTokenRevoke.EXPECT().
		IsTokenOrUserRevoked(gomock.Any(), tokenString, userID.String()).
		Return(false, errors.New("redis connection error"))

	// Act
	validatedUser, err := service.ValidateToken(ctx, tokenString)

	// Assert: should fail closed (return error, not allow through)
	require.Error(t, err)
	assert.Nil(t, validatedUser)

	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrInternalError, apiErr.Code)
}
