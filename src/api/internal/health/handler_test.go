package health

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_Health_WhenAllServicesHealthy_ThenReturnsOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabaseHealthChecker(ctrl)
	mockRedis := NewMockRedisHealthChecker(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(nil)
	mockRedis.EXPECT().Ping(gomock.Any()).Return(redis.NewStatusCmd(context.Background()))
	handler := NewHandler(mockDB, mockRedis)
	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"ok"`)
	assert.Contains(t, string(body), `"postgres":"healthy"`)
	assert.Contains(t, string(body), `"redis":"healthy"`)
}

func Test_Handler_Health_WhenPostgresUnhealthy_ThenReturnsDegraded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabaseHealthChecker(ctrl)
	mockRedis := NewMockRedisHealthChecker(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(errors.New("connection failed"))
	mockRedis.EXPECT().Ping(gomock.Any()).Return(redis.NewStatusCmd(context.Background()))
	handler := NewHandler(mockDB, mockRedis)
	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"degraded"`)
	assert.Contains(t, string(body), `"postgres":"unhealthy"`)
}

func Test_Handler_Health_WhenRedisUnhealthy_ThenReturnsDegraded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabaseHealthChecker(ctrl)
	mockRedis := NewMockRedisHealthChecker(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(nil)
	statusCmd := redis.NewStatusCmd(context.Background())
	statusCmd.SetErr(errors.New("redis unavailable"))
	mockRedis.EXPECT().Ping(gomock.Any()).Return(statusCmd)
	handler := NewHandler(mockDB, mockRedis)
	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"degraded"`)
	assert.Contains(t, string(body), `"redis":"unhealthy"`)
}

func Test_Handler_Health_WhenBothServicesUnhealthy_ThenReturnsDegraded(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabaseHealthChecker(ctrl)
	mockRedis := NewMockRedisHealthChecker(ctrl)
	mockDB.EXPECT().Ping(gomock.Any()).Return(errors.New("db down"))
	statusCmd := redis.NewStatusCmd(context.Background())
	statusCmd.SetErr(errors.New("redis down"))
	mockRedis.EXPECT().Ping(gomock.Any()).Return(statusCmd)
	handler := NewHandler(mockDB, mockRedis)
	app := fiber.New()
	app.Get("/health", handler.Health)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"degraded"`)
	assert.Contains(t, string(body), `"postgres":"unhealthy"`)
	assert.Contains(t, string(body), `"redis":"unhealthy"`)
}

func Test_Handler_Ready_WhenCalled_ThenReturnsReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabaseHealthChecker(ctrl)
	mockRedis := NewMockRedisHealthChecker(ctrl)
	handler := NewHandler(mockDB, mockRedis)
	app := fiber.New()
	app.Get("/ready", handler.Ready)

	req := httptest.NewRequest("GET", "/ready", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"ready"`)
}
