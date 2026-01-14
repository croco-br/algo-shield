package branding

import (
	"context"
	"errors"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_GetBranding_WhenConfigExists_ThenReturnsConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iconURL := "/custom/icon.svg"
	faviconURL := "/custom/favicon.ico"
	expectedConfig := &models.BrandingConfig{
		ID:             1,
		AppName:        "CustomApp",
		IconURL:        &iconURL,
		FaviconURL:     &faviconURL,
		PrimaryColor:   "#FF5733",
		SecondaryColor: "#33FF57",
		HeaderColor:    "#333333",
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any()).Return(expectedConfig, nil)
	service := NewService(mockRepo)

	config, err := service.GetBranding(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}

func Test_Service_GetBranding_WhenConfigDoesNotExist_ThenReturnsDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().Get(gomock.Any()).Return(nil, errors.New("not found"))
	service := NewService(mockRepo)

	config, err := service.GetBranding(context.Background())

	require.NoError(t, err)
	assert.Equal(t, DefaultAppName, config.AppName)
	assert.Equal(t, DefaultPrimaryColor, config.PrimaryColor)
	assert.Equal(t, DefaultSecondaryColor, config.SecondaryColor)
	assert.Equal(t, DefaultHeaderColor, config.HeaderColor)
	assert.NotNil(t, config.IconURL)
	assert.Equal(t, DefaultIconURL, *config.IconURL)
	assert.NotNil(t, config.FaviconURL)
	assert.Equal(t, DefaultFaviconURL, *config.FaviconURL)
}

func Test_Service_UpdateBranding_WhenValidRequest_ThenUpdatesConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	iconURL := "/new/icon.svg"
	req := &UpdateBrandingRequest{
		AppName:        "NewApp",
		IconURL:        &iconURL,
		PrimaryColor:   "#3B82F6",
		SecondaryColor: "#10B981",
		HeaderColor:    "#1e1e1e",
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	service := NewService(mockRepo)

	config, err := service.UpdateBranding(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, req.AppName, config.AppName)
	assert.Equal(t, req.PrimaryColor, config.PrimaryColor)
	assert.Equal(t, req.SecondaryColor, config.SecondaryColor)
	assert.Equal(t, req.HeaderColor, config.HeaderColor)
}

func Test_Service_UpdateBranding_WhenInvalidColorFormat_ThenReturnsError(t *testing.T) {
	testCases := []struct {
		name          string
		request       *UpdateBrandingRequest
		expectedError string
	}{
		{
			name: "invalid primary color",
			request: &UpdateBrandingRequest{
				AppName:        "NewApp",
				PrimaryColor:   "invalid-color",
				SecondaryColor: "#10B981",
				HeaderColor:    "#1e1e1e",
			},
			expectedError: "primary_color",
		},
		{
			name: "invalid secondary color",
			request: &UpdateBrandingRequest{
				AppName:        "NewApp",
				PrimaryColor:   "#3B82F6",
				SecondaryColor: "not-a-color",
				HeaderColor:    "#1e1e1e",
			},
			expectedError: "secondary_color",
		},
		{
			name: "invalid header color",
			request: &UpdateBrandingRequest{
				AppName:        "NewApp",
				PrimaryColor:   "#3B82F6",
				SecondaryColor: "#10B981",
				HeaderColor:    "bad-hex",
			},
			expectedError: "header_color",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepository(ctrl)
			service := NewService(mockRepo)

			config, err := service.UpdateBranding(context.Background(), tc.request)

			assert.Nil(t, config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func Test_Service_UpdateBranding_WhenAppNameTooLong_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}
	req := &UpdateBrandingRequest{
		AppName:        string(longName),
		PrimaryColor:   "#3B82F6",
		SecondaryColor: "#10B981",
		HeaderColor:    "#1e1e1e",
	}
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	config, err := service.UpdateBranding(context.Background(), req)

	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100 characters")
}

func Test_Service_UpdateBranding_WhenAppNameEmpty_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &UpdateBrandingRequest{
		AppName:        "",
		PrimaryColor:   "#3B82F6",
		SecondaryColor: "#10B981",
		HeaderColor:    "#1e1e1e",
	}
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)

	config, err := service.UpdateBranding(context.Background(), req)

	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func Test_Service_UpdateBranding_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &UpdateBrandingRequest{
		AppName:        "NewApp",
		PrimaryColor:   "#3B82F6",
		SecondaryColor: "#10B981",
		HeaderColor:    "#1e1e1e",
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
	service := NewService(mockRepo)

	config, err := service.UpdateBranding(context.Background(), req)

	assert.Nil(t, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update branding configuration")
}

func Test_Service_UpdateBranding_WhenShortHexColor_ThenAccepts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &UpdateBrandingRequest{
		AppName:        "NewApp",
		PrimaryColor:   "#FFF",
		SecondaryColor: "#000",
		HeaderColor:    "#ABC",
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	service := NewService(mockRepo)

	config, err := service.UpdateBranding(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "#FFF", config.PrimaryColor)
	assert.Equal(t, "#000", config.SecondaryColor)
	assert.Equal(t, "#ABC", config.HeaderColor)
}
