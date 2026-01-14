//go:build integration

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Postgres *pgxpool.Pool
	Redis    *redisclient.Client
	Cleanup  func()
}

func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(10*time.Second)),
	)
	require.NoError(t, err)

	postgresConnStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	redisEndpoint, err := redisContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, postgresConnStr)
	require.NoError(t, err)

	redisClient := redisclient.NewClient(&redisclient.Options{
		Addr: redisEndpoint,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	require.NoError(t, err)

	err = redisClient.Ping(ctx).Err()
	require.NoError(t, err)

	runMigrations(t, pool)

	cleanup := func() {
		pool.Close()
		redisClient.Close()
		cleanupCtx := context.Background()
		if err := postgresContainer.Terminate(cleanupCtx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
		if err := redisContainer.Terminate(cleanupCtx); err != nil {
			t.Logf("failed to terminate redis container: %v", err)
		}
	}

	t.Cleanup(cleanup)

	return &TestDB{
		Postgres: pool,
		Redis:    redisClient,
		Cleanup:  cleanup,
	}
}

func runMigrations(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	basePath := "../../../../scripts/migrations"

	// Read all .sql files from the migrations directory
	entries, err := os.ReadDir(basePath)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Collect all .sql files
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			// Skip the migrations.sh script file if it exists
			if entry.Name() != "migrations.sh" {
				migrationFiles = append(migrationFiles, entry.Name())
			}
		}
	}

	// Sort files by name to ensure correct execution order
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		t.Fatalf("No migration files found in %s", basePath)
	}

	// Execute migrations in order
	for _, filename := range migrationFiles {
		filePath := filepath.Join(basePath, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			t.Fatalf("Failed to run migration %s: %v", filename, err)
		}
		t.Logf("Successfully executed migration: %s", filename)
	}
}

func TruncateTables(t *testing.T, pool *pgxpool.Pool, tables []string) {
	t.Helper()
	ctx := context.Background()

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
}
