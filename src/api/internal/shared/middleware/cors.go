package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() fiber.Handler {
	// Enhanced CORS configuration for better browser compatibility, especially Brave
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length, Content-Type",
		MaxAge:           86400,
	})
}

// SecurityHeaders adds security headers that work well with Brave browser
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add headers that help with Brave browser compatibility
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Continue to next handler
		return c.Next()
	}
}
