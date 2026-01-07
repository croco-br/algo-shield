//go:build integration

package permissions_test

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/permissions"
	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_PermissionsUserRepository_GetUserByID_ReturnsUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetUserByID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.Name)
}

func TestIntegration_PermissionsUserRepository_GetUserByID_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetUserByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_PermissionsUserRepository_ListUsers_ReturnsAllUsers(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)
	`, userID1, "user1@example.com", "User 1", "local", true, time.Now(), time.Now(),
		userID2, "user2@example.com", "User 2", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.ListUsers(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestIntegration_PermissionsUserRepository_UpdateUserActive_UpdatesActiveStatus(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	err = repo.UpdateUserActive(ctx, userID, false)

	require.NoError(t, err)

	var active bool
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT active FROM users WHERE id = $1
	`, userID).Scan(&active)
	require.NoError(t, err)
	assert.False(t, active)
}

func TestIntegration_PermissionsUserRepository_CountActiveAdmins_ReturnsCount(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	adminRoleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, adminRoleID, "admin", "Admin role", time.Now(), time.Now())
	require.NoError(t, err)

	userID1 := uuid.New()
	userID2 := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)
	`, userID1, "admin1@example.com", "Admin 1", "local", true, time.Now(), time.Now(),
		userID2, "admin2@example.com", "Admin 2", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3), ($4, $2, $5)
	`, userID1, adminRoleID, time.Now(), userID2, time.Now())
	require.NoError(t, err)

	count, err := repo.CountActiveAdmins(ctx, nil)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2)
}

func TestIntegration_PermissionsUserRepository_CountActiveAdmins_ExcludeUser_ExcludesUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	adminRoleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, adminRoleID, "admin", "Admin role", time.Now(), time.Now())
	require.NoError(t, err)

	userID1 := uuid.New()
	userID2 := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)
	`, userID1, "admin1@example.com", "Admin 1", "local", true, time.Now(), time.Now(),
		userID2, "admin2@example.com", "Admin 2", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3), ($4, $2, $5)
	`, userID1, adminRoleID, time.Now(), userID2, time.Now())
	require.NoError(t, err)

	count, err := repo.CountActiveAdmins(ctx, &userID1)

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestIntegration_PermissionsUserRepository_HasAdminRole_DirectRole_ReturnsTrue(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	adminRoleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, adminRoleID, "admin", "Admin role", time.Now(), time.Now())
	require.NoError(t, err)

	userID := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "admin@example.com", "Admin User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, adminRoleID, time.Now())
	require.NoError(t, err)

	hasRole, err := repo.HasAdminRole(ctx, userID)

	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestIntegration_PermissionsUserRepository_HasAdminRole_GroupRole_ReturnsTrue(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	adminRoleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, adminRoleID, "admin", "Admin role", time.Now(), time.Now())
	require.NoError(t, err)

	groupID := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "admin_group", "Admin group", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO group_roles (group_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, groupID, adminRoleID, time.Now())
	require.NoError(t, err)

	userID := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "user@example.com", "User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_groups (user_id, group_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, groupID, time.Now())
	require.NoError(t, err)

	hasRole, err := repo.HasAdminRole(ctx, userID)

	require.NoError(t, err)
	assert.True(t, hasRole)
}

func TestIntegration_PermissionsUserRepository_HasAdminRole_NoRole_ReturnsFalse(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := permissions.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "user@example.com", "User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	hasRole, err := repo.HasAdminRole(ctx, userID)

	require.NoError(t, err)
	assert.False(t, hasRole)
}
