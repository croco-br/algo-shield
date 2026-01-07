package roles

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_ListRoles_WhenSuccess_ThenReturnsRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "viewer", Description: "Viewer", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().ListRoles(gomock.Any()).Return(expectedRoles, nil)
	service := NewService(mockRepo)

	roles, err := service.ListRoles(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)
}

func Test_Service_ListRoles_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().ListRoles(gomock.Any()).Return(nil, errors.New("database error"))
	service := NewService(mockRepo)

	roles, err := service.ListRoles(context.Background())

	assert.Nil(t, roles)
	assert.Error(t, err)
}

func Test_Service_GetRoleByID_WhenRoleExists_ThenReturnsRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleID := uuid.New()
	expectedRole := &models.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRoleByID(gomock.Any(), roleID).Return(expectedRole, nil)
	service := NewService(mockRepo)

	role, err := service.GetRoleByID(context.Background(), roleID)

	require.NoError(t, err)
	assert.Equal(t, expectedRole, role)
}

func Test_Service_GetRoleByID_WhenRoleNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRoleByID(gomock.Any(), roleID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo)

	role, err := service.GetRoleByID(context.Background(), roleID)

	assert.Nil(t, role)
	assert.Error(t, err)
}

func Test_Service_GetRoleByName_WhenRoleExists_ThenReturnsRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedRole := &models.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRoleByName(gomock.Any(), "admin").Return(expectedRole, nil)
	service := NewService(mockRepo)

	role, err := service.GetRoleByName(context.Background(), "admin")

	require.NoError(t, err)
	assert.Equal(t, expectedRole, role)
}

func Test_Service_GetRoleByName_WhenRoleNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRoleByName(gomock.Any(), "nonexistent").Return(nil, errors.New("not found"))
	service := NewService(mockRepo)

	role, err := service.GetRoleByName(context.Background(), "nonexistent")

	assert.Nil(t, role)
	assert.Error(t, err)
}

func Test_Service_AssignRole_WhenSuccess_ThenReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	roleID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().AssignRole(gomock.Any(), userID, roleID).Return(nil)
	service := NewService(mockRepo)

	err := service.AssignRole(context.Background(), userID, roleID)

	assert.NoError(t, err)
}

func Test_Service_AssignRole_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	roleID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().AssignRole(gomock.Any(), userID, roleID).Return(errors.New("database error"))
	service := NewService(mockRepo)

	err := service.AssignRole(context.Background(), userID, roleID)

	assert.Error(t, err)
}

func Test_Service_RemoveRole_WhenSuccess_ThenReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	roleID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().RemoveRole(gomock.Any(), userID, roleID).Return(nil)
	service := NewService(mockRepo)

	err := service.RemoveRole(context.Background(), userID, roleID)

	assert.NoError(t, err)
}

func Test_Service_RemoveRole_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	roleID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().RemoveRole(gomock.Any(), userID, roleID).Return(errors.New("database error"))
	service := NewService(mockRepo)

	err := service.RemoveRole(context.Background(), userID, roleID)

	assert.Error(t, err)
}

func Test_Service_LoadUserRoles_WhenSuccess_ThenReturnsRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator"},
		{ID: uuid.New(), Name: "editor", Description: "Editor"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserRoles(gomock.Any(), userID).Return(expectedRoles, nil)
	service := NewService(mockRepo)

	roles, err := service.LoadUserRoles(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)
}

func Test_Service_LoadUserRoles_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserRoles(gomock.Any(), userID).Return(nil, errors.New("database error"))
	service := NewService(mockRepo)

	roles, err := service.LoadUserRoles(context.Background(), userID)

	assert.Nil(t, roles)
	assert.Error(t, err)
}

func Test_Service_LoadUserRoles_WhenUserHasNoRoles_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserRoles(gomock.Any(), userID).Return([]models.Role{}, nil)
	service := NewService(mockRepo)

	roles, err := service.LoadUserRoles(context.Background(), userID)

	require.NoError(t, err)
	assert.Empty(t, roles)
}
