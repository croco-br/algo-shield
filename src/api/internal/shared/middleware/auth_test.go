package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/algo-shield/algo-shield/src/api/internal/auth"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AuthMiddleware_WhenNoAuthHeader_ThenReturnsUnauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/test", AuthMiddleware(&auth.Handler{}), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_AuthMiddleware_WhenInvalidAuthHeaderFormat_ThenReturnsUnauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/test", AuthMiddleware(&auth.Handler{}), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_AuthMiddleware_WhenMissingBearerPrefix_ThenReturnsUnauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/test", AuthMiddleware(&auth.Handler{}), func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "token123")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_RequireRole_WhenUserHasRole_ThenAllowsAccess(t *testing.T) {
	middleware := RequireRole("admin")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Admin User",
			Email: "admin@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "admin"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_RequireRole_WhenUserDoesNotHaveRole_ThenReturnsForbidden(t *testing.T) {
	middleware := RequireRole("admin")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Regular User",
			Email: "user@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "user"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func Test_RequireRole_WhenUserNotInContext_ThenReturnsUnauthorized(t *testing.T) {
	middleware := RequireRole("admin")

	app := fiber.New()
	app.Get("/test", middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_RequireAnyRole_WhenUserHasOneOfRoles_ThenAllowsAccess(t *testing.T) {
	middleware := RequireAnyRole("admin", "moderator")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Moderator User",
			Email: "mod@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "moderator"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_RequireAnyRole_WhenUserHasNoneOfRoles_ThenReturnsForbidden(t *testing.T) {
	middleware := RequireAnyRole("admin", "moderator")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Regular User",
			Email: "user@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "user"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func Test_RequireAnyRole_WhenUserNotInContext_ThenReturnsUnauthorized(t *testing.T) {
	middleware := RequireAnyRole("admin", "moderator")

	app := fiber.New()
	app.Get("/test", middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_RequireAnyRole_WhenUserHasMultipleMatchingRoles_ThenAllowsAccess(t *testing.T) {
	middleware := RequireAnyRole("admin", "moderator")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Super User",
			Email: "super@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "admin"},
				{ID: uuid.New(), Name: "moderator"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_RequireRole_WhenUserHasMultipleRoles_ThenChecksCorrectly(t *testing.T) {
	middleware := RequireRole("viewer")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "Multi Role User",
			Email: "multi@example.com",
			Roles: []models.Role{
				{ID: uuid.New(), Name: "admin"},
				{ID: uuid.New(), Name: "viewer"},
				{ID: uuid.New(), Name: "editor"},
			},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_RequireRole_WhenUserHasNoRoles_ThenReturnsForbidden(t *testing.T) {
	middleware := RequireRole("admin")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "No Roles User",
			Email: "noroles@example.com",
			Roles: []models.Role{},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

func Test_RequireAnyRole_WhenUserHasNoRoles_ThenReturnsForbidden(t *testing.T) {
	middleware := RequireAnyRole("admin", "moderator")

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Name:  "No Roles User",
			Email: "noroles@example.com",
			Roles: []models.Role{},
		}
		c.Locals("user", user)
		return c.Next()
	}, middleware, func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}
