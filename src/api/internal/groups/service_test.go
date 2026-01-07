package groups

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

func Test_Service_ListGroups_WhenSuccess_ThenReturnsGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedGroups := []models.Group{
		{ID: uuid.New(), Name: "admins", Description: "Administrators", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New(), Name: "editors", Description: "Editors", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().ListGroups(gomock.Any()).Return(expectedGroups, nil)
	service := NewService(mockRepo)

	groups, err := service.ListGroups(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expectedGroups, groups)
}

func Test_Service_ListGroups_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().ListGroups(gomock.Any()).Return(nil, errors.New("database error"))
	service := NewService(mockRepo)

	groups, err := service.ListGroups(context.Background())

	assert.Nil(t, groups)
	assert.Error(t, err)
}

func Test_Service_GetGroupByID_WhenGroupExists_ThenReturnsGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupID := uuid.New()
	expectedGroup := &models.Group{
		ID:          groupID,
		Name:        "admins",
		Description: "Administrators",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetGroupByID(gomock.Any(), groupID).Return(expectedGroup, nil)
	service := NewService(mockRepo)

	group, err := service.GetGroupByID(context.Background(), groupID)

	require.NoError(t, err)
	assert.Equal(t, expectedGroup, group)
}

func Test_Service_GetGroupByID_WhenGroupNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetGroupByID(gomock.Any(), groupID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo)

	group, err := service.GetGroupByID(context.Background(), groupID)

	assert.Nil(t, group)
	assert.Error(t, err)
}

func Test_Service_GetGroupByName_WhenGroupExists_ThenReturnsGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedGroup := &models.Group{
		ID:          uuid.New(),
		Name:        "admins",
		Description: "Administrators",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetGroupByName(gomock.Any(), "admins").Return(expectedGroup, nil)
	service := NewService(mockRepo)

	group, err := service.GetGroupByName(context.Background(), "admins")

	require.NoError(t, err)
	assert.Equal(t, expectedGroup, group)
}

func Test_Service_GetGroupByName_WhenGroupNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetGroupByName(gomock.Any(), "nonexistent").Return(nil, errors.New("not found"))
	service := NewService(mockRepo)

	group, err := service.GetGroupByName(context.Background(), "nonexistent")

	assert.Nil(t, group)
	assert.Error(t, err)
}

func Test_Service_LoadUserGroups_WhenSuccess_ThenReturnsGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedGroups := []models.Group{
		{ID: uuid.New(), Name: "admins", Description: "Administrators"},
		{ID: uuid.New(), Name: "editors", Description: "Editors"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserGroups(gomock.Any(), userID).Return(expectedGroups, nil)
	service := NewService(mockRepo)

	groups, err := service.LoadUserGroups(context.Background(), userID)

	require.NoError(t, err)
	assert.Equal(t, expectedGroups, groups)
}

func Test_Service_LoadUserGroups_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserGroups(gomock.Any(), userID).Return(nil, errors.New("database error"))
	service := NewService(mockRepo)

	groups, err := service.LoadUserGroups(context.Background(), userID)

	assert.Nil(t, groups)
	assert.Error(t, err)
}

func Test_Service_LoadUserGroups_WhenUserHasNoGroups_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadUserGroups(gomock.Any(), userID).Return([]models.Group{}, nil)
	service := NewService(mockRepo)

	groups, err := service.LoadUserGroups(context.Background(), userID)

	require.NoError(t, err)
	assert.Empty(t, groups)
}

func Test_Service_LoadGroupRoles_WhenSuccess_ThenReturnsRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupID := uuid.New()
	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator"},
		{ID: uuid.New(), Name: "editor", Description: "Editor"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadGroupRoles(gomock.Any(), groupID).Return(expectedRoles, nil)
	service := NewService(mockRepo)

	roles, err := service.LoadGroupRoles(context.Background(), groupID)

	require.NoError(t, err)
	assert.Equal(t, expectedRoles, roles)
}

func Test_Service_LoadGroupRoles_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadGroupRoles(gomock.Any(), groupID).Return(nil, errors.New("database error"))
	service := NewService(mockRepo)

	roles, err := service.LoadGroupRoles(context.Background(), groupID)

	assert.Nil(t, roles)
	assert.Error(t, err)
}

func Test_Service_LoadGroupRoles_WhenGroupHasNoRoles_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	groupID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().LoadGroupRoles(gomock.Any(), groupID).Return([]models.Role{}, nil)
	service := NewService(mockRepo)

	roles, err := service.LoadGroupRoles(context.Background(), groupID)

	require.NoError(t, err)
	assert.Empty(t, roles)
}
