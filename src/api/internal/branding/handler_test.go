package branding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_GetBranding_WhenSuccess_ThenReturnsConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iconURL := "/icon.svg"
	expectedConfig := &models.BrandingConfig{
		ID:             1,
		AppName:        "TestApp",
		IconURL:        &iconURL,
		PrimaryColor:   "#3B82F6",
		SecondaryColor: "#10B981",
		HeaderColor:    "#1e1e1e",
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetBranding(gomock.Any()).Return(expectedConfig, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/branding", handler.GetBranding)

	req := httptest.NewRequest("GET", "/branding", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result models.BrandingConfig
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.Equal(t, expectedConfig.AppName, result.AppName)
	assert.Equal(t, expectedConfig.PrimaryColor, result.PrimaryColor)
}

func Test_Handler_GetBranding_WhenServiceFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	mockService.EXPECT().GetBranding(gomock.Any()).Return(nil, errors.New("service error"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/branding", handler.GetBranding)

	req := httptest.NewRequest("GET", "/branding", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Failed to fetch branding configuration")
}

func Test_Handler_UpdateBranding_WhenValidRequest_ThenUpdatesConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iconURL := "/new/icon.svg"
	requestBody := UpdateBrandingRequest{
		AppName:        "UpdatedApp",
		IconURL:        &iconURL,
		PrimaryColor:   "#FF5733",
		SecondaryColor: "#33FF57",
		HeaderColor:    "#333333",
	}
	updatedConfig := &models.BrandingConfig{
		ID:             1,
		AppName:        requestBody.AppName,
		IconURL:        requestBody.IconURL,
		PrimaryColor:   requestBody.PrimaryColor,
		SecondaryColor: requestBody.SecondaryColor,
		HeaderColor:    requestBody.HeaderColor,
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().UpdateBranding(gomock.Any(), gomock.Any()).Return(updatedConfig, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/branding", handler.UpdateBranding)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/branding", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result models.BrandingConfig
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.Equal(t, updatedConfig.AppName, result.AppName)
}

func Test_Handler_UpdateBranding_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/branding", handler.UpdateBranding)

	req := httptest.NewRequest("PUT", "/branding", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid request body")
}

func Test_Handler_UpdateBranding_WhenValidationFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestBody := UpdateBrandingRequest{
		AppName:        "",
		PrimaryColor:   "#FF5733",
		SecondaryColor: "#33FF57",
		HeaderColor:    "#333333",
	}
	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/branding", handler.UpdateBranding)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/branding", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_UpdateBranding_WhenServiceFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	requestBody := UpdateBrandingRequest{
		AppName:        "UpdatedApp",
		PrimaryColor:   "#FF5733",
		SecondaryColor: "#33FF57",
		HeaderColor:    "#333333",
	}
	mockService := NewMockService(ctrl)
	mockService.EXPECT().UpdateBranding(gomock.Any(), gomock.Any()).Return(nil, errors.New("update failed"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/branding", handler.UpdateBranding)
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("PUT", "/branding", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "update failed")
}
