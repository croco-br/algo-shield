package tokenrevoke

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Service handles JWT token revocation using Redis blacklist
type Service struct {
	redis *redis.Client
}

// NewService creates a new token revocation service
func NewService(redisClient *redis.Client) *Service {
	return &Service{
		redis: redisClient,
	}
}

// RevokeToken adds a token to the blacklist
// The token is hashed before storing for security
// TTL is set to match the token's expiration time
func (s *Service) RevokeToken(ctx context.Context, token string, expiresAt time.Time) error {
	tokenHash := hashToken(token)
	key := fmt.Sprintf("blacklist:%s", tokenHash)

	// Calculate TTL based on token expiration
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Store in Redis with TTL
	err := s.redis.SetEx(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// IsTokenRevoked checks if a token is in the blacklist
func (s *Service) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	tokenHash := hashToken(token)
	key := fmt.Sprintf("blacklist:%s", tokenHash)

	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		// If Redis is unavailable, log error but allow token through (availability > security for reads)
		// In production, you might want to handle this differently based on your requirements
		return false, fmt.Errorf("failed to check token revocation status: %w", err)
	}

	return exists > 0, nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
// This is useful when a user changes password or is deactivated
func (s *Service) RevokeAllUserTokens(ctx context.Context, userID string, tokenExpiry time.Duration) error {
	// We use a user-specific blacklist key that will block all tokens for this user
	key := fmt.Sprintf("blacklist:user:%s", userID)

	// Store with TTL equal to the maximum token expiry time
	err := s.redis.SetEx(ctx, key, "1", tokenExpiry).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke all user tokens: %w", err)
	}

	return nil
}

// IsUserTokensRevoked checks if all tokens for a user are revoked
func (s *Service) IsUserTokensRevoked(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("blacklist:user:%s", userID)

	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check user tokens revocation status: %w", err)
	}

	return exists > 0, nil
}

// hashToken creates a SHA-256 hash of the token for storage
// This prevents storing the actual token in Redis
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Health checks if the Redis connection is healthy
func (s *Service) Health(ctx context.Context) error {
	return s.redis.Ping(ctx).Err()
}
