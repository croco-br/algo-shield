package system

import (
	"bytes"
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

func Test_Handler_GetSyntheticMode_WhenSuccess_ThenReturnsConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	response := &SyntheticModeResponse{
		Enabled:   true,
		UpdatedAt: time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetSyntheticMode(gomock.Any()).Return(response, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/system/mode", handler.GetSyntheticMode)

	req := httptest.NewRequest("GET", "/system/mode", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result SyntheticModeResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.True(t, result.Enabled)
}

func Test_Handler_GetSyntheticMode_WhenServiceFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetSyntheticMode(gomock.Any()).Return(nil, errors.New("service error"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/system/mode", handler.GetSyntheticMode)

	req := httptest.NewRequest("GET", "/system/mode", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to get synthetic mode")
}

func Test_Handler_GetSyntheticMode_WhenDisabled_ThenReturnsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	response := &SyntheticModeResponse{
		Enabled:   false,
		UpdatedAt: time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetSyntheticMode(gomock.Any()).Return(response, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/system/mode", handler.GetSyntheticMode)

	req := httptest.NewRequest("GET", "/system/mode", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result SyntheticModeResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.False(t, result.Enabled)
}

func Test_Handler_SetSyntheticMode_WhenValidRequest_ThenUpdatesConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestBody := UpdateSyntheticModeRequest{Enabled: true}
	response := &SyntheticModeResponse{
		Enabled:   true,
		UpdatedAt: time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().SetSyntheticMode(gomock.Any(), true).Return(response, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/system/mode", handler.SetSyntheticMode)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/system/mode", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result SyntheticModeResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.True(t, result.Enabled)
}

func Test_Handler_SetSyntheticMode_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/system/mode", handler.SetSyntheticMode)

	req := httptest.NewRequest("PUT", "/system/mode", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid request body")
}

func Test_Handler_SetSyntheticMode_WhenServiceFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestBody := UpdateSyntheticModeRequest{Enabled: true}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().SetSyntheticMode(gomock.Any(), true).Return(nil, errors.New("service error"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/system/mode", handler.SetSyntheticMode)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/system/mode", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to update synthetic mode")
}

func Test_Handler_SetSyntheticMode_WhenDisabling_ThenUpdatesConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestBody := UpdateSyntheticModeRequest{Enabled: false}
	response := &SyntheticModeResponse{
		Enabled:   false,
		UpdatedAt: time.Now(),
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().SetSyntheticMode(gomock.Any(), false).Return(response, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/system/mode", handler.SetSyntheticMode)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/system/mode", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result SyntheticModeResponse
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.False(t, result.Enabled)
}
