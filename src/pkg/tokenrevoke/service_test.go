package tokenrevoke

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test_Service_RevokeToken_WhenValidToken_ThenStoresInRedis tests token revocation
func Test_Service_RevokeToken_WhenValidToken_ThenStoresInRedis(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	token := "test.jwt.token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Expect SetEx to be called with correct parameters
	mockRedis.EXPECT().
		SetEx(ctx, gomock.Any(), "1", gomock.Any()).
		DoAndReturn(func(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
			// Verify key format
			assert.Contains(t, key, "blacklist:")
			// Verify TTL is positive
			assert.Greater(t, ttl, time.Duration(0))
			return redis.NewStatusCmd(ctx)
		})

	// Act
	err := service.RevokeToken(ctx, token, expiresAt)

	// Assert
	require.NoError(t, err)
}

// Test_Service_RevokeToken_WhenExpiredToken_ThenDoesNotStore tests expired token handling
func Test_Service_RevokeToken_WhenExpiredToken_ThenDoesNotStore(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	token := "test.jwt.token"
	expiresAt := time.Now().Add(-1 * time.Hour) // Already expired

	// No Redis call should be made for expired token

	// Act
	err := service.RevokeToken(ctx, token, expiresAt)

	// Assert
	require.NoError(t, err)
}

// Test_Service_IsTokenRevoked_WhenTokenBlacklisted_ThenReturnsTrue tests revoked token check
func Test_Service_IsTokenRevoked_WhenTokenBlacklisted_ThenReturnsTrue(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	token := "test.jwt.token"

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, keys ...string) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(1) // Token exists in blacklist
			return cmd
		})

	// Act
	isRevoked, err := service.IsTokenRevoked(ctx, token)

	// Assert
	require.NoError(t, err)
	assert.True(t, isRevoked)
}

// Test_Service_IsTokenRevoked_WhenTokenNotBlacklisted_ThenReturnsFalse tests non-revoked token
func Test_Service_IsTokenRevoked_WhenTokenNotBlacklisted_ThenReturnsFalse(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	token := "test.jwt.token"

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, keys ...string) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(0) // Token does not exist in blacklist
			return cmd
		})

	// Act
	isRevoked, err := service.IsTokenRevoked(ctx, token)

	// Assert
	require.NoError(t, err)
	assert.False(t, isRevoked)
}

// Test_Service_RevokeAllUserTokens_WhenValidUserID_ThenStoresInRedis tests user-level revocation
func Test_Service_RevokeAllUserTokens_WhenValidUserID_ThenStoresInRedis(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	userID := "user-123"
	tokenExpiry := 24 * time.Hour

	mockRedis.EXPECT().
		SetEx(ctx, gomock.Any(), "1", tokenExpiry).
		DoAndReturn(func(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
			// Verify key format includes user ID
			assert.Contains(t, key, "blacklist:user:")
			assert.Contains(t, key, userID)
			return redis.NewStatusCmd(ctx)
		})

	// Act
	err := service.RevokeAllUserTokens(ctx, userID, tokenExpiry)

	// Assert
	require.NoError(t, err)
}

// Test_Service_IsUserTokensRevoked_WhenUserTokensRevoked_ThenReturnsTrue tests user revocation check
func Test_Service_IsUserTokensRevoked_WhenUserTokensRevoked_ThenReturnsTrue(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	userID := "user-123"

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, keys ...string) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(1) // User tokens are revoked
			return cmd
		})

	// Act
	isRevoked, err := service.IsUserTokensRevoked(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.True(t, isRevoked)
}

// Test_Service_IsUserTokensRevoked_WhenUserTokensNotRevoked_ThenReturnsFalse tests non-revoked user
func Test_Service_IsUserTokensRevoked_WhenUserTokensNotRevoked_ThenReturnsFalse(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()
	userID := "user-123"

	mockRedis.EXPECT().
		Exists(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, keys ...string) *redis.IntCmd {
			cmd := redis.NewIntCmd(ctx)
			cmd.SetVal(0) // User tokens are not revoked
			return cmd
		})

	// Act
	isRevoked, err := service.IsUserTokensRevoked(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.False(t, isRevoked)
}

// Test_Service_hashToken_WhenSameToken_ThenReturnsSameHash tests deterministic hashing
func Test_Service_hashToken_WhenSameToken_ThenReturnsSameHash(t *testing.T) {
	// Arrange
	token := "test.jwt.token"

	// Act
	hash1 := hashToken(token)
	hash2 := hashToken(token)

	// Assert
	assert.Equal(t, hash1, hash2, "Same token should produce same hash")
	assert.NotEmpty(t, hash1)
	assert.Len(t, hash1, 64, "SHA-256 hash should be 64 hex characters")
}

// Test_Service_hashToken_WhenDifferentTokens_ThenReturnsDifferentHashes tests hash uniqueness
func Test_Service_hashToken_WhenDifferentTokens_ThenReturnsDifferentHashes(t *testing.T) {
	// Arrange
	token1 := "test.jwt.token1"
	token2 := "test.jwt.token2"

	// Act
	hash1 := hashToken(token1)
	hash2 := hashToken(token2)

	// Assert
	assert.NotEqual(t, hash1, hash2, "Different tokens should produce different hashes")
}

// Test_Service_Health_WhenRedisHealthy_ThenReturnsNoError tests health check
func Test_Service_Health_WhenRedisHealthy_ThenReturnsNoError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := NewMockRedisClient(ctrl)
	service := &Service{redis: mockRedis}

	ctx := context.Background()

	mockRedis.EXPECT().
		Ping(ctx).
		DoAndReturn(func(ctx context.Context) *redis.StatusCmd {
			cmd := redis.NewStatusCmd(ctx)
			cmd.SetVal("PONG")
			return cmd
		})

	// Act
	err := service.Health(ctx)

	// Assert
	require.NoError(t, err)
}
