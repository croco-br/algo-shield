//go:build integration

package roles_test

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RolesRepository_ListRoles_ReturnsAllRoles(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.ListRoles(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.GreaterOrEqual(t, len(result), 3)
}

func TestIntegration_RolesRepository_GetRoleByID_ReturnsRole(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	roleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "test_role", "Test role description", time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetRoleByID(ctx, roleID)

	require.NoError(t, err)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "test_role", result.Name)
	assert.Equal(t, "Test role description", result.Description)
}

func TestIntegration_RolesRepository_GetRoleByID_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetRoleByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_RolesRepository_GetRoleByName_ReturnsRole(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	roleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "unique_role", "Unique role", time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetRoleByName(ctx, "unique_role")

	require.NoError(t, err)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "unique_role", result.Name)
}

func TestIntegration_RolesRepository_GetRoleByName_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.GetRoleByName(ctx, "non_existent_role")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_RolesRepository_AssignRole_AssignsRoleToUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "test_role", "Test role", time.Now(), time.Now())
	require.NoError(t, err)

	err = repo.AssignRole(ctx, userID, roleID)

	require.NoError(t, err)

	var count int
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2
	`, userID, roleID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestIntegration_RolesRepository_AssignRole_Duplicate_NoError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "test_role", "Test role", time.Now(), time.Now())
	require.NoError(t, err)

	err = repo.AssignRole(ctx, userID, roleID)
	require.NoError(t, err)

	err = repo.AssignRole(ctx, userID, roleID)

	require.NoError(t, err)
}

func TestIntegration_RolesRepository_RemoveRole_RemovesRoleFromUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "test_role", "Test role", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, roleID, time.Now())
	require.NoError(t, err)

	err = repo.RemoveRole(ctx, userID, roleID)

	require.NoError(t, err)

	var count int
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2
	`, userID, roleID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestIntegration_RolesRepository_LoadUserRoles_ReturnsUserRoles(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID1 := uuid.New()
	roleID2 := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)
	`, roleID1, "role1", "Role 1", time.Now(), time.Now(),
		roleID2, "role2", "Role 2", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3), ($1, $4, $5)
	`, userID, roleID1, time.Now(), roleID2, time.Now())
	require.NoError(t, err)

	result, err := repo.LoadUserRoles(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntegration_RolesRepository_GetRoleIDByName_WithTransaction_ReturnsRoleID(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	roleID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "tx_role", "Transaction role", time.Now(), time.Now())
	require.NoError(t, err)

	tx, err := testDB.Postgres.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	result, err := repo.GetRoleIDByName(ctx, tx, "tx_role")

	require.NoError(t, err)
	assert.Equal(t, roleID, result)
}

func TestIntegration_RolesRepository_AssignRoleToUser_WithTransaction_AssignsRole(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "tx_role", "Transaction role", time.Now(), time.Now())
	require.NoError(t, err)

	tx, err := testDB.Postgres.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	err = repo.AssignRoleToUser(ctx, tx, userID, roleID)

	require.NoError(t, err)

	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2
	`, userID, roleID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestIntegration_RolesRepository_GetRoleIDByName_NotFound_ReturnsNilUUID(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	tx, err := testDB.Postgres.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	result, err := repo.GetRoleIDByName(ctx, tx, "non_existent_role")

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
}

func TestIntegration_RolesRepository_AssignRoleToUser_Duplicate_NoError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "tx_role", "Transaction role", time.Now(), time.Now())
	require.NoError(t, err)

	tx, err := testDB.Postgres.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	err = repo.AssignRoleToUser(ctx, tx, userID, roleID)
	require.NoError(t, err)

	err = repo.AssignRoleToUser(ctx, tx, userID, roleID)

	require.NoError(t, err)

	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2
	`, userID, roleID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestIntegration_RolesRepository_RemoveRole_NonExistent_NoError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "test_role", "Test role", time.Now(), time.Now())
	require.NoError(t, err)

	err = repo.RemoveRole(ctx, userID, roleID)

	require.NoError(t, err)
}

func TestIntegration_RolesRepository_LoadUserRoles_WithGroupRoles_ReturnsAllRoles(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID1 := uuid.New()
	roleID2 := uuid.New()
	groupID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)
	`, roleID1, "direct_role", "Direct role", time.Now(), time.Now(),
		roleID2, "group_role", "Group role", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "test_group", "Test group", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, roleID1, time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_groups (user_id, group_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, groupID, time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO group_roles (group_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, groupID, roleID2, time.Now())
	require.NoError(t, err)

	result, err := repo.LoadUserRoles(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntegration_RolesRepository_LoadUserRoles_DuplicateRoles_ReturnsUnique(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := roles.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	roleID := uuid.New()
	groupID := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, roleID, "duplicate_role", "Duplicate role", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "test_group", "Test group", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, roleID, time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_groups (user_id, group_id, assigned_at)
		VALUES ($1, $2, $3)
	`, userID, groupID, time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO group_roles (group_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
	`, groupID, roleID, time.Now())
	require.NoError(t, err)

	result, err := repo.LoadUserRoles(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, roleID, result[0].ID)
}
