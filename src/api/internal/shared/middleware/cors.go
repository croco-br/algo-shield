package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS creates a CORS middleware with the specified allowed origins
// allowedOrigins: comma-separated list of allowed origins (e.g., "https://example.com,https://app.example.com")
// Use "*" only for development/testing (NOT production)
func CORS(allowedOrigins string) fiber.Handler {
	// SECURITY: In production, allowedOrigins must be specific domains (validated in config)
	// Example production value: "https://yourdomain.com,https://app.yourdomain.com"
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Synthetic-Mode",
		AllowCredentials: false, // Keep false for security (JWT in Authorization header)
		ExposeHeaders:    "Content-Length, Content-Type",
		MaxAge:           86400, // 24 hours cache for preflight requests
	})
}

// SecurityHeaders adds comprehensive security headers following OWASP best practices
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking attacks
		c.Set("X-Frame-Options", "DENY")

		// Control referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Enable XSS protection (legacy browsers that don't support CSP)
		c.Set("X-XSS-Protection", "1; mode=block")

		// SECURITY: HSTS header for HTTPS enforcement (max-age=1 year)
		// Only set when using HTTPS to avoid browser warnings
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// SECURITY: Content Security Policy to prevent XSS attacks
		// This is a strict policy - adjust based on your application needs
		csp := "default-src 'self'; " +
			"script-src 'self'; " +
			"style-src 'self' 'unsafe-inline'; " + // unsafe-inline needed for some UI frameworks
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Set("Content-Security-Policy", csp)

		// Prevent cross-domain policy abuse
		c.Set("X-Permitted-Cross-Domain-Policies", "none")

		// Continue to next handler
		return c.Next()
	}
}
