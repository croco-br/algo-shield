package middleware

import (
	"bytes"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logger_WhenRequestSucceeds_ThenLogsRequest(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	middleware := Logger()

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "200")
}

func Test_Logger_WhenRequestFails_ThenLogsError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	middleware := Logger()

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).SendString("error")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "500")
}

func Test_Logger_WhenMultipleRequests_ThenLogsAll(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	middleware := Logger()

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test1", func(c *fiber.Ctx) error {
		return c.SendString("success1")
	})
	app.Get("/test2", func(c *fiber.Ctx) error {
		return c.SendString("success2")
	})

	req1 := httptest.NewRequest("GET", "/test1", nil)
	req2 := httptest.NewRequest("GET", "/test2", nil)

	resp1, err1 := app.Test(req1)
	resp2, err2 := app.Test(req2)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "/test1")
	assert.Contains(t, logOutput, "/test2")
}

func Test_Logger_WhenPostRequest_ThenLogsMethod(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	middleware := Logger()

	app := fiber.New()
	app.Use(middleware)
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("created")
	})

	req := httptest.NewRequest("POST", "/test", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "POST")
	assert.Contains(t, logOutput, "/test")
}
