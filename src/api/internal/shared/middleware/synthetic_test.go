package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_SyntheticModeMiddleware_WhenHeaderIsTrue_ThenSetsSyntheticMode(t *testing.T) {
	app := fiber.New()
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(SyntheticModeHeader, "true")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func Test_SyntheticModeMiddleware_WhenHeaderIsFalse_ThenDoesNotSetSyntheticMode(t *testing.T) {
	app := fiber.New()
	var capturedMode bool
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		capturedMode = IsSyntheticMode(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(SyntheticModeHeader, "false")
	_, err := app.Test(req)

	assert.NoError(t, err)
	assert.False(t, capturedMode)
}

func Test_SyntheticModeMiddleware_WhenHeaderMissing_ThenDefaultsToFalse(t *testing.T) {
	app := fiber.New()
	var capturedMode bool
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		capturedMode = IsSyntheticMode(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)

	assert.NoError(t, err)
	assert.False(t, capturedMode)
}

func Test_IsSyntheticMode_WhenModeSetToTrue_ThenReturnsTrue(t *testing.T) {
	app := fiber.New()
	var result bool
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		result = IsSyntheticMode(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(SyntheticModeHeader, "true")
	_, err := app.Test(req)

	assert.NoError(t, err)
	assert.True(t, result)
}

func Test_IsSyntheticMode_WhenModeNotSet_ThenReturnsFalse(t *testing.T) {
	app := fiber.New()
	var result bool
	app.Get("/test", func(c *fiber.Ctx) error {
		result = IsSyntheticMode(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	_, err := app.Test(req)

	assert.NoError(t, err)
	assert.False(t, result)
}

func Test_GetTableName_WhenSyntheticModeTrue_ThenAddsSyntheticPrefix(t *testing.T) {
	tableName := GetTableName(true, "transactions")

	assert.Equal(t, "synthetic_transactions", tableName)
}

func Test_GetTableName_WhenSyntheticModeFalse_ThenReturnsBaseTable(t *testing.T) {
	tableName := GetTableName(false, "transactions")

	assert.Equal(t, "transactions", tableName)
}

func Test_GetTableName_WithDifferentBaseTable_ThenAppliesPrefix(t *testing.T) {
	tableName := GetTableName(true, "users")

	assert.Equal(t, "synthetic_users", tableName)
}

func Test_SyntheticModeMiddleware_WhenHeaderCaseMismatch_ThenHandlesCorrectly(t *testing.T) {
	app := fiber.New()
	var capturedMode bool
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		capturedMode = IsSyntheticMode(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(SyntheticModeHeader, "TRUE")
	_, err := app.Test(req)

	assert.NoError(t, err)
	assert.False(t, capturedMode)
}

func Test_SyntheticModeMiddleware_WhenMultipleRequests_ThenIsolatesContext(t *testing.T) {
	app := fiber.New()
	var modes []bool
	app.Use(SyntheticModeMiddleware())
	app.Get("/test", func(c *fiber.Ctx) error {
		modes = append(modes, IsSyntheticMode(c))
		return c.SendString("ok")
	})

	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.Header.Set(SyntheticModeHeader, "true")
	_, _ = app.Test(req1)

	req2 := httptest.NewRequest("GET", "/test", nil)
	_, _ = app.Test(req2)

	assert.Len(t, modes, 2)
	assert.True(t, modes[0])
	assert.False(t, modes[1])
}
