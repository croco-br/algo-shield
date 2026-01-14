package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	MaxRequests int           // Maximum requests per window
	Window      time.Duration // Time window for rate limiting
	KeyPrefix   string        // Redis key prefix (e.g., "ratelimit:login:")
}

// RateLimiter creates a rate limiting middleware using Redis for distributed rate limiting
// Uses sliding window algorithm for accurate rate limiting
func RateLimiter(redisClient *redis.Client, config RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client identifier (IP address)
		// In production, you might want to use a combination of IP + user agent
		clientIP := c.IP()
		key := fmt.Sprintf("%s%s", config.KeyPrefix, clientIP)

		ctx := context.Background()

		// Sliding window counter implementation
		now := time.Now()
		windowStart := now.Add(-config.Window)

		// Use Redis pipeline for efficiency
		pipe := redisClient.Pipeline()

		// Remove old entries outside the window
		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart.UnixNano()))

		// Count current requests in window
		countCmd := pipe.ZCard(ctx, key)

		// Execute pipeline
		_, err := pipe.Exec(ctx)
		if err != nil {
			// On Redis error, allow the request (fail open for availability)
			// Log the error in production
			return c.Next()
		}

		currentCount := countCmd.Val()

		// Check if rate limit exceeded
		if currentCount >= int64(config.MaxRequests) {
			// Rate limit exceeded
			c.Set("X-RateLimit-Limit", fmt.Sprint(config.MaxRequests))
			c.Set("X-RateLimit-Remaining", "0")
			c.Set("X-RateLimit-Reset", fmt.Sprint(now.Add(config.Window).Unix()))

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Please try again in %v", config.Window),
			})
		}

		// Add current request to the window
		score := now.UnixNano()
		member := redis.Z{
			Score:  float64(score),
			Member: fmt.Sprintf("%d", score),
		}
		_, err = redisClient.ZAdd(ctx, key, member).Result()
		if err != nil {
			// On Redis error, allow the request (fail open)
			return c.Next()
		}

		// Set TTL on the key to auto-cleanup
		redisClient.Expire(ctx, key, config.Window*2)

		// Set rate limit headers
		remaining := config.MaxRequests - int(currentCount) - 1
		if remaining < 0 {
			remaining = 0
		}
		c.Set("X-RateLimit-Limit", fmt.Sprint(config.MaxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprint(remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprint(now.Add(config.Window).Unix()))

		return c.Next()
	}
}

// Default rate limit configurations for different endpoints
var (
	// LoginRateLimit: 5 attempts per minute per IP
	LoginRateLimit = RateLimitConfig{
		MaxRequests: 5,
		Window:      1 * time.Minute,
		KeyPrefix:   "ratelimit:login:",
	}

	// RegisterRateLimit: 3 attempts per hour per IP
	RegisterRateLimit = RateLimitConfig{
		MaxRequests: 3,
		Window:      1 * time.Hour,
		KeyPrefix:   "ratelimit:register:",
	}

	// APIRateLimit: 100 requests per minute per IP for general API endpoints
	APIRateLimit = RateLimitConfig{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		KeyPrefix:   "ratelimit:api:",
	}
)
