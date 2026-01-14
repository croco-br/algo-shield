package csrf

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// TokenLength is the length of the CSRF token in bytes (32 bytes = 256 bits)
	TokenLength = 32

	// TokenTTL is the time-to-live for CSRF tokens in Redis (24 hours)
	TokenTTL = 24 * time.Hour

	// HeaderName is the name of the header that contains the CSRF token
	HeaderName = "X-CSRF-Token"
)

// GenerateToken generates a cryptographically secure random CSRF token
func GenerateToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// StoreToken stores a CSRF token in Redis with TTL
func StoreToken(ctx context.Context, redis *redis.Client, userID string, token string) error {
	key := fmt.Sprintf("csrf:%s", userID)
	return redis.Set(ctx, key, token, TokenTTL).Err()
}

// ValidateToken validates a CSRF token against the stored token in Redis
func ValidateToken(ctx context.Context, redis *redis.Client, userID string, token string) bool {
	if token == "" {
		return false
	}

	key := fmt.Sprintf("csrf:%s", userID)
	storedToken, err := redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	return storedToken == token
}

// DeleteToken deletes a CSRF token from Redis (used on logout)
func DeleteToken(ctx context.Context, redis *redis.Client, userID string) error {
	key := fmt.Sprintf("csrf:%s", userID)
	return redis.Del(ctx, key).Err()
}
