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
	password := "password123"

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
	user, token, err := service.LoginUser(ctx, email, password)

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
	password := "password123"

	mockUserService.EXPECT().
		GetUserByEmailWithPassword(ctx, email).
		Return(nil, errors.New("user not found"))

	// Act
	user, token, err := service.LoginUser(ctx, email, password)

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
	correctPassword := "password123"
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
	user, token, err := service.LoginUser(ctx, email, wrongPassword)

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
	password := "password123"

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
	user, token, err := service.LoginUser(ctx, email, password)

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

	// Set up mock expectations before calling ValidateToken
	mockTokenRevoke.EXPECT().
		IsTokenRevoked(gomock.Any(), tokenString).
		Return(false, nil).
		Times(1)

	mockTokenRevoke.EXPECT().
		IsUserTokensRevoked(gomock.Any(), userID.String()).
		Return(false, nil).
		Times(1)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(user, nil).
		Times(1)

	// Act
	validatedUser, err := service.ValidateToken(tokenString)

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

	// Act
	user, err := service.ValidateToken(invalidToken)

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

	// Mock token as revoked
	mockTokenRevoke.EXPECT().
		IsTokenRevoked(gomock.Any(), tokenString).
		Return(true, nil)

	// Act
	validatedUser, err := service.ValidateToken(tokenString)

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

	// Token itself is not revoked
	mockTokenRevoke.EXPECT().
		IsTokenRevoked(gomock.Any(), tokenString).
		Return(false, nil)

	// But all user tokens are revoked (e.g., password change)
	mockTokenRevoke.EXPECT().
		IsUserTokensRevoked(gomock.Any(), userID.String()).
		Return(true, nil)

	// Act
	validatedUser, err := service.ValidateToken(tokenString)

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
	password := "password123"

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
	user, token, err := service.RegisterUser(ctx, email, name, password)

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
	password := "password123"

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
	user, token, err := service.RegisterUser(ctx, email, name, password)

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
	password := "password123"

	mockUserService.EXPECT().
		GetUserByEmail(ctx, email).
		Return(nil, errors.New("not found"))

	mockUserService.EXPECT().
		CreateUser(ctx, email, name, gomock.Any()).
		Return(nil, errors.New("database error"))

	user, token, err := service.RegisterUser(ctx, email, name, password)

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
	password := "password123"

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

	user, token, err := service.LoginUser(ctx, email, password)

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

	user, err := service.ValidateToken(invalidToken)

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

	mockTokenRevoke.EXPECT().
		IsTokenRevoked(gomock.Any(), tokenString).
		Return(false, nil)

	mockTokenRevoke.EXPECT().
		IsUserTokensRevoked(gomock.Any(), userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("user not found"))

	validatedUser, err := service.ValidateToken(tokenString)

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

	inactiveUser := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: false,
	}

	mockTokenRevoke.EXPECT().
		IsTokenRevoked(gomock.Any(), tokenString).
		Return(false, nil)

	mockTokenRevoke.EXPECT().
		IsUserTokensRevoked(gomock.Any(), userID.String()).
		Return(false, nil)

	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(inactiveUser, nil)

	validatedUser, err := service.ValidateToken(tokenString)

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
