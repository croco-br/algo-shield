package permissions

import (
	"context"
	"errors"
	"testing"

	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test_Service_UpdateUserActive_WhenActivatingUser_ThenSucceeds tests user activation
func Test_Service_UpdateUserActive_WhenActivatingUser_ThenSucceeds(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	currentUserID := uuid.New()
	targetUserID := uuid.New()

	mockUserRepo.EXPECT().
		UpdateUserActive(ctx, targetUserID, true).
		Return(nil)

	// Act
	err := service.UpdateUserActive(ctx, currentUserID, targetUserID, true)

	// Assert
	require.NoError(t, err)
}

// Test_Service_UpdateUserActive_WhenDeactivatingSelf_ThenReturnsCannotDeactivateSelfError tests admin protection
func Test_Service_UpdateUserActive_WhenDeactivatingSelf_ThenReturnsCannotDeactivateSelfError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	userID := uuid.New()

	// Act
	err := service.UpdateUserActive(ctx, userID, userID, false)

	// Assert
	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrCannotDeactivateSelf, apiErr.Code)
}

// Test_Service_UpdateUserActive_WhenDeactivatingLastAdmin_ThenReturnsCannotDeactivateLastAdminError tests last admin protection
func Test_Service_UpdateUserActive_WhenDeactivatingLastAdmin_ThenReturnsCannotDeactivateLastAdminError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	currentUserID := uuid.New()
	targetUserID := uuid.New()

	// Target user is an admin
	mockUserRepo.EXPECT().
		HasAdminRole(ctx, targetUserID).
		Return(true, nil)

	// No other active admins
	mockUserRepo.EXPECT().
		CountActiveAdmins(ctx, &targetUserID).
		Return(0, nil)

	// Act
	err := service.UpdateUserActive(ctx, currentUserID, targetUserID, false)

	// Assert
	require.Error(t, err)
	apiErr, ok := err.(*apierrors.APIError)
	require.True(t, ok, "Expected APIError")
	assert.Equal(t, apierrors.ErrCannotDeactivateLastAdmin, apiErr.Code)
}

// Test_Service_UpdateUserActive_WhenDeactivatingAdminWithOthersActive_ThenSucceeds tests valid admin deactivation
func Test_Service_UpdateUserActive_WhenDeactivatingAdminWithOthersActive_ThenSucceeds(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	currentUserID := uuid.New()
	targetUserID := uuid.New()

	// Target user is an admin
	mockUserRepo.EXPECT().
		HasAdminRole(ctx, targetUserID).
		Return(true, nil)

	// There are other active admins
	mockUserRepo.EXPECT().
		CountActiveAdmins(ctx, &targetUserID).
		Return(2, nil) // 2 other active admins

	mockUserRepo.EXPECT().
		UpdateUserActive(ctx, targetUserID, false).
		Return(nil)

	// Act
	err := service.UpdateUserActive(ctx, currentUserID, targetUserID, false)

	// Assert
	require.NoError(t, err)
}

// Test_Service_UpdateUserActive_WhenDeactivatingNonAdmin_ThenSucceeds tests non-admin deactivation
func Test_Service_UpdateUserActive_WhenDeactivatingNonAdmin_ThenSucceeds(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	currentUserID := uuid.New()
	targetUserID := uuid.New()

	// Target user is not an admin
	mockUserRepo.EXPECT().
		HasAdminRole(ctx, targetUserID).
		Return(false, nil)

	mockUserRepo.EXPECT().
		UpdateUserActive(ctx, targetUserID, false).
		Return(nil)

	// Act
	err := service.UpdateUserActive(ctx, currentUserID, targetUserID, false)

	// Assert
	require.NoError(t, err)
}

// Test_Service_GetUserByID_WhenUserExists_ThenReturnsUserWithRolesAndGroups tests user retrieval
func Test_Service_GetUserByID_WhenUserExists_ThenReturnsUserWithRolesAndGroups(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	userID := uuid.New()

	user := &models.User{
		ID:     userID,
		Email:  "test@example.com",
		Name:   "Test User",
		Active: true,
	}

	roles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator"},
	}

	groups := []models.Group{
		{ID: uuid.New(), Name: "developers", Description: "Development team"},
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(user, nil)

	mockRoleService.EXPECT().
		LoadUserRoles(ctx, userID).
		Return(roles, nil)

	mockGroupService.EXPECT().
		LoadUserGroups(ctx, userID).
		Return(groups, nil)

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, roles, result.Roles)
	assert.Equal(t, groups, result.Groups)
}

// Test_Service_GetUserByID_WhenUserNotFound_ThenReturnsError tests user not found
func Test_Service_GetUserByID_WhenUserNotFound_ThenReturnsError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.EXPECT().
		GetUserByID(ctx, userID).
		Return(nil, errors.New("user not found"))

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
}

// Test_Service_ListUsers_WhenUsersExist_ThenReturnsUsersWithRolesAndGroups tests user listing
func Test_Service_ListUsers_WhenUsersExist_ThenReturnsUsersWithRolesAndGroups(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockRoleService := NewMockService(ctrl)
	mockGroupService := NewMockGroupService(ctrl)

	service := NewService(mockUserRepo, mockRoleService, mockGroupService)

	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()

	users := []models.User{
		{ID: userID1, Email: "user1@example.com", Name: "User 1", Active: true},
		{ID: userID2, Email: "user2@example.com", Name: "User 2", Active: true},
	}

	roles1 := []models.Role{{ID: uuid.New(), Name: "admin", Description: "Administrator"}}
	roles2 := []models.Role{{ID: uuid.New(), Name: "user", Description: "Regular user"}}

	groups1 := []models.Group{{ID: uuid.New(), Name: "group1", Description: "Group 1"}}
	groups2 := []models.Group{{ID: uuid.New(), Name: "group2", Description: "Group 2"}}

	mockUserRepo.EXPECT().
		ListUsers(ctx).
		Return(users, nil)

	// Load roles and groups for each user
	mockRoleService.EXPECT().
		LoadUserRoles(ctx, userID1).
		Return(roles1, nil)

	mockGroupService.EXPECT().
		LoadUserGroups(ctx, userID1).
		Return(groups1, nil)

	mockRoleService.EXPECT().
		LoadUserRoles(ctx, userID2).
		Return(roles2, nil)

	mockGroupService.EXPECT().
		LoadUserGroups(ctx, userID2).
		Return(groups2, nil)

	// Act
	result, err := service.ListUsers(ctx)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, roles1, result[0].Roles)
	assert.Equal(t, groups1, result[0].Groups)
	assert.Equal(t, roles2, result[1].Roles)
	assert.Equal(t, groups2, result[1].Groups)
}
