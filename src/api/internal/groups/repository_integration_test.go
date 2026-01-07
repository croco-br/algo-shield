//go:build integration

package groups_test

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_GroupsRepository_ListGroups_ReturnsAllGroups(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.ListGroups(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestIntegration_GroupsRepository_GetGroupByID_ReturnsGroup(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	groupID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "test_group", "Test group description", time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetGroupByID(ctx, groupID)

	require.NoError(t, err)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "test_group", result.Name)
	assert.Equal(t, "Test group description", result.Description)
}

func TestIntegration_GroupsRepository_GetGroupByID_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetGroupByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_GroupsRepository_GetGroupByName_ReturnsGroup(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	groupID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "unique_group", "Unique group", time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetGroupByName(ctx, "unique_group")

	require.NoError(t, err)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "unique_group", result.Name)
}

func TestIntegration_GroupsRepository_GetGroupByName_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.GetGroupByName(ctx, "non_existent_group")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_GroupsRepository_LoadUserGroups_ReturnsUserGroups(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	groupID1 := uuid.New()
	groupID2 := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)
	`, groupID1, "group1", "Group 1", time.Now(), time.Now(),
		groupID2, "group2", "Group 2", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO user_groups (user_id, group_id, assigned_at)
		VALUES ($1, $2, $3), ($1, $4, $5)
	`, userID, groupID1, time.Now(), groupID2, time.Now())
	require.NoError(t, err)

	result, err := repo.LoadUserGroups(ctx, userID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntegration_GroupsRepository_LoadUserGroups_NoGroups_ReturnsEmpty(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.LoadUserGroups(ctx, userID)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestIntegration_GroupsRepository_LoadGroupRoles_ReturnsGroupRoles(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	groupID := uuid.New()
	roleID1 := uuid.New()
	roleID2 := uuid.New()

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "test_group", "Test group", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)
	`, roleID1, "role1", "Role 1", time.Now(), time.Now(),
		roleID2, "role2", "Role 2", time.Now(), time.Now())
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO group_roles (group_id, role_id, assigned_at)
		VALUES ($1, $2, $3), ($1, $4, $5)
	`, groupID, roleID1, time.Now(), roleID2, time.Now())
	require.NoError(t, err)

	result, err := repo.LoadGroupRoles(ctx, groupID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIntegration_GroupsRepository_LoadGroupRoles_NoRoles_ReturnsEmpty(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := groups.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	groupID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, groupID, "test_group", "Test group", time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.LoadGroupRoles(ctx, groupID)

	require.NoError(t, err)
	assert.Empty(t, result)
}
