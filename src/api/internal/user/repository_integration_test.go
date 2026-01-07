//go:build integration

package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/algo-shield/algo-shield/src/api/internal/user"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestIntegration_UserRepository_GetUserByEmail_ReturnsUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	email := "test@example.com"
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, email, "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetUserByEmail(ctx, email, false)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, "Test User", result.Name)
	assert.Nil(t, result.PasswordHash)
}

func TestIntegration_UserRepository_GetUserByEmail_WithPassword_ReturnsPasswordHash(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	email := "test@example.com"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, password_hash, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, userID, email, "Test User", string(passwordHash), "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetUserByEmail(ctx, email, true)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.NotNil(t, result.PasswordHash)
	assert.Equal(t, string(passwordHash), *result.PasswordHash)
}

func TestIntegration_UserRepository_GetUserByEmail_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.GetUserByEmail(ctx, "nonexistent@example.com", false)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_UserRepository_GetUserByID_ReturnsUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
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
}

func TestIntegration_UserRepository_GetUserByID_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetUserByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_UserRepository_CreateUser_StoresUser(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	passwordHash := "hashed_password"
	now := time.Now()
	testUser := &models.User{
		ID:           userID,
		Email:        "newuser@example.com",
		Name:         "New User",
		PasswordHash: &passwordHash,
		AuthType:     models.AuthTypeLocal,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := repo.CreateUser(ctx, testUser)

	require.NoError(t, err)

	var storedEmail string
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT email FROM users WHERE id = $1
	`, userID).Scan(&storedEmail)
	require.NoError(t, err)
	assert.Equal(t, "newuser@example.com", storedEmail)
}

func TestIntegration_UserRepository_CreateUser_DuplicateEmail_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	email := "duplicate@example.com"
	userID1 := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID1, email, "First User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	userID2 := uuid.New()
	passwordHash := "hashed_password"
	now := time.Now()
	testUser := &models.User{
		ID:           userID2,
		Email:        email,
		Name:         "Second User",
		PasswordHash: &passwordHash,
		AuthType:     models.AuthTypeLocal,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = repo.CreateUser(ctx, testUser)

	assert.Error(t, err)
}

func TestIntegration_UserRepository_CreateUserWithTx_StoresUserInTransaction(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	tx, err := testDB.Postgres.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	userID := uuid.New()
	passwordHash := "hashed_password"
	now := time.Now()
	testUser := &models.User{
		ID:           userID,
		Email:        "txuser@example.com",
		Name:         "TX User",
		PasswordHash: &passwordHash,
		AuthType:     models.AuthTypeLocal,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = repo.CreateUserWithTx(ctx, tx, testUser)

	require.NoError(t, err)

	var storedEmail string
	err = tx.QueryRow(ctx, `
		SELECT email FROM users WHERE id = $1
	`, userID).Scan(&storedEmail)
	require.NoError(t, err)
	assert.Equal(t, "txuser@example.com", storedEmail)
}

func TestIntegration_UserRepository_UpdateLastLogin_UpdatesLastLoginAt(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now())
	require.NoError(t, err)

	lastLoginAt := time.Now()
	err = repo.UpdateLastLogin(ctx, userID, &lastLoginAt)

	require.NoError(t, err)

	var storedLastLogin *time.Time
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT last_login_at FROM users WHERE id = $1
	`, userID).Scan(&storedLastLogin)
	require.NoError(t, err)
	assert.NotNil(t, storedLastLogin)
}

func TestIntegration_UserRepository_UpdateLastLogin_NilTime_ClearsLastLogin(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := user.NewPostgresUserRepository(testDB.Postgres)
	ctx := context.Background()

	userID := uuid.New()
	lastLoginAt := time.Now()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO users (id, email, name, auth_type, active, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, userID, "test@example.com", "Test User", "local", true, time.Now(), time.Now(), lastLoginAt)
	require.NoError(t, err)

	err = repo.UpdateLastLogin(ctx, userID, nil)

	require.NoError(t, err)

	var storedLastLogin *time.Time
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT last_login_at FROM users WHERE id = $1
	`, userID).Scan(&storedLastLogin)
	require.NoError(t, err)
}
