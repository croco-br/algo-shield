package system

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_GetSyntheticMode_WhenSuccess_ThenReturnsConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &SyntheticModeConfig{Enabled: true}
	updatedAt := time.Now()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(config, updatedAt, nil)
	service := NewService(mockRepo)

	result, err := service.GetSyntheticMode(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, config.Enabled, result.Enabled)
	assert.Equal(t, updatedAt, result.UpdatedAt)
}

func Test_Service_GetSyntheticMode_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(nil, time.Time{}, errors.New("database error"))
	service := NewService(mockRepo)

	result, err := service.GetSyntheticMode(context.Background())

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func Test_Service_GetSyntheticMode_WhenDisabled_ThenReturnsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &SyntheticModeConfig{Enabled: false}
	updatedAt := time.Now()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(config, updatedAt, nil)
	service := NewService(mockRepo)

	result, err := service.GetSyntheticMode(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Enabled)
}

func Test_Service_SetSyntheticMode_WhenEnabling_ThenReturnsUpdatedConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &SyntheticModeConfig{Enabled: true}
	updatedAt := time.Now()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().SetSyntheticMode(gomock.Any(), true).Return(nil)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(config, updatedAt, nil)
	service := NewService(mockRepo)

	result, err := service.SetSyntheticMode(context.Background(), true)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Enabled)
	assert.Equal(t, updatedAt, result.UpdatedAt)
}

func Test_Service_SetSyntheticMode_WhenDisabling_ThenReturnsUpdatedConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &SyntheticModeConfig{Enabled: false}
	updatedAt := time.Now()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().SetSyntheticMode(gomock.Any(), false).Return(nil)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(config, updatedAt, nil)
	service := NewService(mockRepo)

	result, err := service.SetSyntheticMode(context.Background(), false)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Enabled)
}

func Test_Service_SetSyntheticMode_WhenSetFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().SetSyntheticMode(gomock.Any(), true).Return(errors.New("database error"))
	service := NewService(mockRepo)

	result, err := service.SetSyntheticMode(context.Background(), true)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func Test_Service_SetSyntheticMode_WhenGetAfterSetFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().SetSyntheticMode(gomock.Any(), true).Return(nil)
	mockRepo.EXPECT().GetSyntheticMode(gomock.Any()).Return(nil, time.Time{}, errors.New("get error"))
	service := NewService(mockRepo)

	result, err := service.SetSyntheticMode(context.Background(), true)

	assert.Nil(t, result)
	assert.Error(t, err)
}
