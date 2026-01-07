//go:build integration

package branding_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/branding"
	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_BrandingRepository_Get_ReturnsConfig(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	result, err := repo.Get(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.NotEmpty(t, result.AppName)
}

func TestIntegration_BrandingRepository_Get_UsesCache(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	result1, err := repo.Get(ctx)
	require.NoError(t, err)

	cachedConfig, err := testDB.Redis.Get(ctx, "branding:config").Result()
	require.NoError(t, err)

	var cached models.BrandingConfig
	err = json.Unmarshal([]byte(cachedConfig), &cached)
	require.NoError(t, err)

	assert.Equal(t, result1.ID, cached.ID)
	assert.Equal(t, result1.AppName, cached.AppName)
}

func TestIntegration_BrandingRepository_Get_WithoutRedis_ReturnsFromDB(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, nil)
	ctx := context.Background()

	result, err := repo.Get(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
}

func TestIntegration_BrandingRepository_Update_UpdatesConfig(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	iconURL := "https://example.com/icon.png"
	faviconURL := "https://example.com/favicon.ico"
	config := &models.BrandingConfig{
		ID:             1,
		AppName:        "Updated App",
		IconURL:        &iconURL,
		FaviconURL:     &faviconURL,
		PrimaryColor:   "#FF0000",
		SecondaryColor: "#00FF00",
		HeaderColor:    "#0000FF",
		UpdatedAt:      time.Now(),
	}

	err := repo.Update(ctx, config)

	require.NoError(t, err)

	result, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, "Updated App", result.AppName)
	assert.Equal(t, "#FF0000", result.PrimaryColor)
	assert.Equal(t, "#00FF00", result.SecondaryColor)
	assert.Equal(t, "#0000FF", result.HeaderColor)
	if result.IconURL != nil {
		assert.Equal(t, iconURL, *result.IconURL)
	}
	if result.FaviconURL != nil {
		assert.Equal(t, faviconURL, *result.FaviconURL)
	}
}

func TestIntegration_BrandingRepository_Update_InvalidatesCache(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	_, err := repo.Get(ctx)
	require.NoError(t, err)

	exists, err := testDB.Redis.Exists(ctx, "branding:config").Result()
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	config := &models.BrandingConfig{
		ID:             1,
		AppName:        "Cache Test",
		PrimaryColor:   "#FF0000",
		SecondaryColor: "#00FF00",
		HeaderColor:    "#0000FF",
		UpdatedAt:      time.Now(),
	}

	err = repo.Update(ctx, config)
	require.NoError(t, err)

	exists, err = testDB.Redis.Exists(ctx, "branding:config").Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

func TestIntegration_BrandingRepository_Update_WithoutRedis_UpdatesDB(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, nil)
	ctx := context.Background()

	config := &models.BrandingConfig{
		ID:             1,
		AppName:        "No Redis App",
		PrimaryColor:   "#FF0000",
		SecondaryColor: "#00FF00",
		HeaderColor:    "#0000FF",
		UpdatedAt:      time.Now(),
	}

	err := repo.Update(ctx, config)

	require.NoError(t, err)

	result, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, "No Redis App", result.AppName)
}

func TestIntegration_BrandingRepository_Get_AfterUpdate_ReturnsUpdatedConfig(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := branding.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	original, err := repo.Get(ctx)
	require.NoError(t, err)

	config := &models.BrandingConfig{
		ID:             1,
		AppName:        "Final App Name",
		PrimaryColor:   "#123456",
		SecondaryColor: "#789ABC",
		HeaderColor:    "#DEF012",
		UpdatedAt:      time.Now(),
	}

	err = repo.Update(ctx, config)
	require.NoError(t, err)

	updated, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.NotEqual(t, original.AppName, updated.AppName)
	assert.Equal(t, "Final App Name", updated.AppName)
	assert.Equal(t, "#123456", updated.PrimaryColor)
}
