package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CORS_WhenCalled_ThenSetsCORSHeaders(t *testing.T) {
	// Test with wildcard for development
	middleware := CORS("*")

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func Test_CORS_WhenOptionsRequest_ThenHandlesPreflight(t *testing.T) {
	middleware := CORS("*")

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func Test_SecurityHeaders_WhenCalled_ThenSetsSecurityHeaders(t *testing.T) {
	middleware := SecurityHeaders()

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify all security headers are set
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'")
	assert.Equal(t, "none", resp.Header.Get("X-Permitted-Cross-Domain-Policies"))

	// HSTS should NOT be set when using HTTP (test environment)
	assert.Empty(t, resp.Header.Get("Strict-Transport-Security"))
}

func Test_SecurityHeaders_WhenChainedWithOtherMiddleware_ThenWorksCorrectly(t *testing.T) {
	securityMiddleware := SecurityHeaders()
	corsMiddleware := CORS("*")

	app := fiber.New()
	app.Use(corsMiddleware)
	app.Use(securityMiddleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
}

// Test_CORS_WithSpecificOrigins tests CORS with specific allowed origins (production scenario)
func Test_CORS_WithSpecificOrigins_ThenAllowsOnlySpecifiedOrigins(t *testing.T) {
	middleware := CORS("https://example.com,https://app.example.com")

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	// Fiber CORS middleware returns the requested origin if it's in the allowed list
	assert.Equal(t, "https://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
}
