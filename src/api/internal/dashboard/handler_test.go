package dashboard

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_GetMetrics_WhenSuccess_ThenReturnsMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedMetrics := &DashboardMetrics{
		StatusDistribution: []StatusCount{{Status: "approved", Count: 10}},
		Temporal24h:        []TemporalCount{{Bucket: time.Now(), Count: 5}},
		Temporal7d:         []TemporalCount{{Bucket: time.Now(), Count: 50}},
		Temporal30d:        []TemporalCount{{Bucket: time.Now(), Count: 200}},
		TotalCount:         100,
		CachedAt:           time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetMetrics(gomock.Any()).Return(expectedMetrics, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/dashboard/metrics", handler.GetMetrics)

	req := httptest.NewRequest("GET", "/dashboard/metrics", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.NotNil(t, result["data"])
	assert.NotNil(t, result["response_time_ms"])
	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(100), data["total_count"])
}

func Test_Handler_GetMetrics_WhenServiceFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetMetrics(gomock.Any()).Return(nil, errors.New("service error"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/dashboard/metrics", handler.GetMetrics)

	req := httptest.NewRequest("GET", "/dashboard/metrics", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to fetch dashboard metrics")
}

func Test_Handler_GetMetrics_WhenSyntheticModeHeader_ThenHandlesSyntheticMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedMetrics := &DashboardMetrics{
		StatusDistribution: []StatusCount{{Status: "pending", Count: 50}},
		Temporal24h:        []TemporalCount{},
		Temporal7d:         []TemporalCount{},
		Temporal30d:        []TemporalCount{},
		TotalCount:         50,
		CachedAt:           time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetMetrics(gomock.Any()).Return(expectedMetrics, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/dashboard/metrics", handler.GetMetrics)

	req := httptest.NewRequest("GET", "/dashboard/metrics", nil)
	req.Header.Set("X-Synthetic-Mode", "true")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.NotNil(t, result["data"])
}

func Test_Handler_GetMetrics_WhenEmptyMetrics_ThenReturnsEmptyArrays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedMetrics := &DashboardMetrics{
		StatusDistribution: []StatusCount{},
		Temporal24h:        []TemporalCount{},
		Temporal7d:         []TemporalCount{},
		Temporal30d:        []TemporalCount{},
		TotalCount:         0,
		CachedAt:           time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetMetrics(gomock.Any()).Return(expectedMetrics, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/dashboard/metrics", handler.GetMetrics)

	req := httptest.NewRequest("GET", "/dashboard/metrics", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total_count"])
}
