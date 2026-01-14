package middleware

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/csrf"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// CSRFConfig holds CSRF protection configuration
type CSRFConfig struct {
	Redis         *redis.Client
	ExcludedPaths []string // Paths to exclude from CSRF validation (e.g., /api/v1/auth/login)
}

// CSRFProtection creates a CSRF protection middleware
// Uses double-submit cookie pattern: token in header + token in Redis keyed by user ID
func CSRFProtection(config CSRFConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip CSRF validation for safe methods (GET, HEAD, OPTIONS)
		method := c.Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return c.Next()
		}

		// Skip CSRF validation for excluded paths
		path := c.Path()
		for _, excluded := range config.ExcludedPaths {
			if path == excluded {
				return c.Next()
			}
		}

		// Get user ID from context (set by auth middleware)
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			// If no user ID, it means request is not authenticated
			// Auth middleware should have already rejected it
			return c.Next()
		}

		// Get CSRF token from header
		csrfToken := c.Get(csrf.HeaderName)
		if csrfToken == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "CSRF token missing",
				"message": "CSRF token is required for this request",
			})
		}

		// Validate CSRF token
		ctx := context.Background()
		if !csrf.ValidateToken(ctx, config.Redis, userID, csrfToken) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Invalid CSRF token",
				"message": "CSRF token is invalid or expired",
			})
		}

		return c.Next()
	}
}
