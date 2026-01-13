package middleware

import (
	"github.com/gofiber/fiber/v2"
)

const (
	// SyntheticModeKey is the key used to store synthetic mode in Fiber context
	SyntheticModeKey = "synthetic_mode"
	// SyntheticModeHeader is the HTTP header used to pass synthetic mode
	SyntheticModeHeader = "X-Synthetic-Mode"
)

// SyntheticModeMiddleware extracts the synthetic mode from the request header
// and stores it in the Fiber context for use by handlers and repositories
func SyntheticModeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		syntheticMode := c.Get(SyntheticModeHeader) == "true"
		c.Locals(SyntheticModeKey, syntheticMode)
		return c.Next()
	}
}

// IsSyntheticMode returns whether synthetic mode is enabled in the request context
func IsSyntheticMode(c *fiber.Ctx) bool {
	if mode, ok := c.Locals(SyntheticModeKey).(bool); ok {
		return mode
	}
	return false
}

// GetTableName returns the appropriate table name based on synthetic mode
func GetTableName(syntheticMode bool, baseTable string) string {
	if syntheticMode {
		return "synthetic_" + baseTable
	}
	return baseTable
}
