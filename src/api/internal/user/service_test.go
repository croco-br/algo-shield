package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_NewService_WhenCalled_ThenReturnsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	assert.NotNil(t, service)
	assert.Equal(t, userRepo, service.userRepo)
	assert.Equal(t, roleRepo, service.roleRepo)
	assert.Equal(t, txManager, service.txManager)
	assert.Equal(t, roleService, service.roleService)
	assert.Equal(t, groupService, service.groupService)
}

func Test_Service_GetUserByEmail_WhenUserExists_ThenReturnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "test@example.com"
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: email,
		Name:  "Test User",
	}

	userRepo.EXPECT().GetUserByEmail(ctx, email, false).Return(expectedUser, nil)
	roleService.EXPECT().LoadUserRoles(ctx, expectedUser.ID).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, expectedUser.ID).Return([]models.Group{}, nil)

	user, err := service.GetUserByEmail(ctx, email)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func Test_Service_GetUserByEmail_WhenUserNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "notfound@example.com"

	userRepo.EXPECT().GetUserByEmail(ctx, email, false).Return(nil, sql.ErrNoRows)

	user, err := service.GetUserByEmail(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)
}

func Test_Service_GetUserByEmailWithPassword_WhenUserExists_ThenReturnsUserWithPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "test@example.com"
	passwordHash := "hashed_password"
	expectedUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         "Test User",
		PasswordHash: &passwordHash,
	}

	userRepo.EXPECT().GetUserByEmail(ctx, email, true).Return(expectedUser, nil)
	roleService.EXPECT().LoadUserRoles(ctx, expectedUser.ID).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, expectedUser.ID).Return([]models.Group{}, nil)

	user, err := service.GetUserByEmailWithPassword(ctx, email)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.NotNil(t, user.PasswordHash)
	assert.Equal(t, passwordHash, *user.PasswordHash)
}

func Test_Service_GetUserByID_WhenUserExists_ThenReturnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	userRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	roleService.EXPECT().LoadUserRoles(ctx, userID).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, userID).Return([]models.Group{}, nil)

	user, err := service.GetUserByID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
}

func Test_Service_GetUserByID_WhenLoadRolesFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	userRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	roleService.EXPECT().LoadUserRoles(ctx, userID).Return(nil, errors.New("database error"))

	user, err := service.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Service_GetUserByID_WhenLoadGroupsFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	userRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)
	roleService.EXPECT().LoadUserRoles(ctx, userID).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, userID).Return(nil, errors.New("database error"))

	user, err := service.GetUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Service_CreateUser_WhenValidData_ThenCreatesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"
	viewerRoleID := uuid.New()

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(sql.ErrTxDone)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, tx pgx.Tx, user *models.User) error {
			assert.Equal(t, email, user.Email)
			assert.Equal(t, name, user.Name)
			assert.NotEqual(t, uuid.Nil, user.ID)
			return nil
		},
	)
	roleRepo.EXPECT().GetRoleIDByName(ctx, mockTx, "viewer").Return(viewerRoleID, nil)
	roleRepo.EXPECT().AssignRoleToUser(ctx, mockTx, gomock.Any(), viewerRoleID).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)
	roleService.EXPECT().LoadUserRoles(ctx, gomock.Any()).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, gomock.Any()).Return([]models.Group{}, nil)

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
}

func Test_Service_CreateUser_WhenTransactionBeginFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"

	txManager.EXPECT().Begin(ctx).Return(nil, errors.New("connection error"))

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Service_CreateUser_WhenCreateUserFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).Return(errors.New("duplicate key"))

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Service_CreateUser_WhenViewerRoleNotFound_ThenCreatesUserWithoutRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(sql.ErrTxDone)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).Return(nil)
	roleRepo.EXPECT().GetRoleIDByName(ctx, mockTx, "viewer").Return(uuid.Nil, sql.ErrNoRows)
	mockTx.EXPECT().Commit(ctx).Return(nil)
	roleService.EXPECT().LoadUserRoles(ctx, gomock.Any()).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, gomock.Any()).Return([]models.Group{}, nil)

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	require.NoError(t, err)
	assert.NotNil(t, user)
}

func Test_Service_CreateUser_WhenViewerRoleIDIsNil_ThenCreatesUserWithoutRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(sql.ErrTxDone)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).Return(nil)
	roleRepo.EXPECT().GetRoleIDByName(ctx, mockTx, "viewer").Return(uuid.Nil, nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)
	roleService.EXPECT().LoadUserRoles(ctx, gomock.Any()).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, gomock.Any()).Return([]models.Group{}, nil)

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	require.NoError(t, err)
	assert.NotNil(t, user)
}

func Test_Service_CreateUser_WhenAssignRoleFails_ThenCreatesUserWithoutRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"
	viewerRoleID := uuid.New()

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(sql.ErrTxDone)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).Return(nil)
	roleRepo.EXPECT().GetRoleIDByName(ctx, mockTx, "viewer").Return(viewerRoleID, nil)
	roleRepo.EXPECT().AssignRoleToUser(ctx, mockTx, gomock.Any(), viewerRoleID).Return(errors.New("constraint violation"))
	mockTx.EXPECT().Commit(ctx).Return(nil)
	roleService.EXPECT().LoadUserRoles(ctx, gomock.Any()).Return([]models.Role{}, nil)
	groupService.EXPECT().LoadUserGroups(ctx, gomock.Any()).Return([]models.Group{}, nil)

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	require.NoError(t, err)
	assert.NotNil(t, user)
}

func Test_Service_CreateUser_WhenCommitFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)
	mockTx := NewMockTx(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	email := "newuser@example.com"
	name := "New User"
	passwordHash := "hashed_password"
	viewerRoleID := uuid.New()

	txManager.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	userRepo.EXPECT().CreateUserWithTx(ctx, mockTx, gomock.Any()).Return(nil)
	roleRepo.EXPECT().GetRoleIDByName(ctx, mockTx, "viewer").Return(viewerRoleID, nil)
	roleRepo.EXPECT().AssignRoleToUser(ctx, mockTx, gomock.Any(), viewerRoleID).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(errors.New("commit failed"))

	user, err := service.CreateUser(ctx, email, name, passwordHash)

	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Service_UpdateLastLogin_WhenCalled_ThenUpdatesLastLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	userID := uuid.New()
	lastLoginAt := time.Now()

	userRepo.EXPECT().UpdateLastLogin(ctx, userID, &lastLoginAt).Return(nil)

	err := service.UpdateLastLogin(ctx, userID, &lastLoginAt)

	require.NoError(t, err)
}

func Test_Service_UpdateLastLogin_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := NewMockUserRepository(ctrl)
	roleRepo := NewMockRoleRepository(ctrl)
	txManager := NewMockTransactionManager(ctrl)
	roleService := NewMockRoleService(ctrl)
	groupService := NewMockGroupService(ctrl)

	service := NewService(userRepo, roleRepo, txManager, roleService, groupService)

	ctx := context.Background()
	userID := uuid.New()
	lastLoginAt := time.Now()

	userRepo.EXPECT().UpdateLastLogin(ctx, userID, &lastLoginAt).Return(errors.New("database error"))

	err := service.UpdateLastLogin(ctx, userID, &lastLoginAt)

	assert.Error(t, err)
}
