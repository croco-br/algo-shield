//go:build integration

package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PostgresRepository_GetStatusDistribution_WhenNoTransactions_ThenReturnsEmptyArray(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	repo := NewPostgresRepository(pool)

	result, err := repo.GetStatusDistribution(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func Test_PostgresRepository_GetStatusDistribution_WhenTransactionsExist_ThenReturnsDistribution(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	testhelpers.InsertTransaction(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "approved", nil, 15, []string{}, `{"amount": 200}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "rejected", nil, 20, []string{"rule1"}, `{"amount": 300}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "in_review", nil, 25, []string{}, `{"amount": 400}`)
	repo := NewPostgresRepository(pool)

	result, err := repo.GetStatusDistribution(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	var approvedCount, rejectedCount, inReviewCount int64
	for _, sc := range result {
		switch sc.Status {
		case "approved":
			approvedCount = sc.Count
		case "rejected":
			rejectedCount = sc.Count
		case "in_review":
			inReviewCount = sc.Count
		}
	}
	assert.Equal(t, int64(2), approvedCount)
	assert.Equal(t, int64(1), rejectedCount)
	assert.Equal(t, int64(1), inReviewCount)
}

func Test_PostgresRepository_GetStatusDistribution_WhenSyntheticMode_ThenUsessyntheticTable(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	testhelpers.TruncateTable(t, pool, "transactions_synthetic")
	testhelpers.InsertTransaction(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`)
	testhelpers.InsertSyntheticTransaction(t, pool, uuid.New(), "pending", nil, 0, []string{}, `{"amount": 200}`)
	testhelpers.InsertSyntheticTransaction(t, pool, uuid.New(), "pending", nil, 0, []string{}, `{"amount": 300}`)
	repo := NewPostgresRepository(pool)
	ctx := shared.WithSyntheticMode(context.Background(), true)

	result, err := repo.GetStatusDistribution(ctx)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "pending", result[0].Status)
	assert.Equal(t, int64(2), result[0].Count)
}

func Test_PostgresRepository_GetTemporalCounts_When24hPeriod_ThenReturnsHourlyCounts(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	now := time.Now()
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`, now.Add(-1*time.Hour))
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 15, []string{}, `{"amount": 200}`, now.Add(-1*time.Hour))
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "rejected", nil, 20, []string{"rule1"}, `{"amount": 300}`, now.Add(-2*time.Hour))
	repo := NewPostgresRepository(pool)

	result, err := repo.GetTemporalCounts(context.Background(), "24h")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
}

func Test_PostgresRepository_GetTemporalCounts_When7dPeriod_ThenReturnsDailyCounts(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	now := time.Now()
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`, now.Add(-24*time.Hour))
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 15, []string{}, `{"amount": 200}`, now.Add(-48*time.Hour))
	repo := NewPostgresRepository(pool)

	result, err := repo.GetTemporalCounts(context.Background(), "7d")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
}

func Test_PostgresRepository_GetTemporalCounts_When30dPeriod_ThenReturnsDailyCounts(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	now := time.Now()
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`, now.Add(-5*24*time.Hour))
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 15, []string{}, `{"amount": 200}`, now.Add(-10*24*time.Hour))
	repo := NewPostgresRepository(pool)

	result, err := repo.GetTemporalCounts(context.Background(), "30d")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
}

func Test_PostgresRepository_GetTemporalCounts_WhenInvalidPeriod_ThenDefaultsTo24h(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	now := time.Now()
	testhelpers.InsertTransactionWithTime(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`, now.Add(-1*time.Hour))
	repo := NewPostgresRepository(pool)

	result, err := repo.GetTemporalCounts(context.Background(), "invalid")

	require.NoError(t, err)
	assert.NotNil(t, result)
}

func Test_PostgresRepository_GetTotalCount_WhenNoTransactions_ThenReturnsZero(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	repo := NewPostgresRepository(pool)

	count, err := repo.GetTotalCount(context.Background())

	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func Test_PostgresRepository_GetTotalCount_WhenTransactionsExist_ThenReturnsCount(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	testhelpers.InsertTransaction(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "rejected", nil, 20, []string{"rule1"}, `{"amount": 200}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "in_review", nil, 15, []string{}, `{"amount": 300}`)
	repo := NewPostgresRepository(pool)

	count, err := repo.GetTotalCount(context.Background())

	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func Test_PostgresRepository_GetTotalCount_WhenSyntheticMode_ThenUsesSyntheticTable(t *testing.T) {
	pool := testhelpers.SetupTestDB(t)
	defer pool.Close()
	testhelpers.TruncateTable(t, pool, "transactions")
	testhelpers.TruncateTable(t, pool, "transactions_synthetic")
	testhelpers.InsertTransaction(t, pool, uuid.New(), "approved", nil, 10, []string{}, `{"amount": 100}`)
	testhelpers.InsertTransaction(t, pool, uuid.New(), "rejected", nil, 20, []string{"rule1"}, `{"amount": 200}`)
	testhelpers.InsertSyntheticTransaction(t, pool, uuid.New(), "pending", nil, 0, []string{}, `{"amount": 300}`)
	repo := NewPostgresRepository(pool)
	ctx := shared.WithSyntheticMode(context.Background(), true)

	count, err := repo.GetTotalCount(ctx)

	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
