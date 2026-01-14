//go:build integration

package transactions_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/algo-shield/algo-shield/src/api/internal/transactions"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TransactionsRepository_GetTransaction_ReturnsTransaction(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	transactionID := uuid.New()
	externalID := "ext-123"
	metadata := map[string]any{"key": "value"}
	metadataJSON, err := json.Marshal(metadata)
	require.NoError(t, err)

	matchedRules := []string{"rule1", "rule2"}
	matchedRulesJSON, err := json.Marshal(matchedRules)
	require.NoError(t, err)

	metadata["external_id"] = externalID
	metadata["amount"] = 100.50
	metadata["currency"] = "USD"
	metadata["origin"] = "account1"
	metadata["destination"] = "account2"
	metadata["type"] = "transfer"
	metadataJSON, err = json.Marshal(metadata)
	require.NoError(t, err)

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO transactions (id, schema_id, status, processing_time, matched_rules, metadata, created_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, transactionID, nil, "approved", 150, matchedRulesJSON, metadataJSON, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetTransaction(ctx, transactionID)

	require.NoError(t, err)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, models.StatusApproved, result.Status)
	assert.Equal(t, int64(150), result.ProcessingTime)
	assert.Len(t, result.MatchedRules, 2)
	assert.Equal(t, externalID, result.Metadata["external_id"])
	assert.Equal(t, 100.50, result.Metadata["amount"])
}

func TestIntegration_TransactionsRepository_GetTransaction_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetTransaction(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_TransactionsRepository_ListTransactions_ReturnsTransactions(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	transactionID1 := uuid.New()
	transactionID2 := uuid.New()
	matchedRulesJSON, _ := json.Marshal([]string{})
	metadataJSON, _ := json.Marshal(map[string]any{})

	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO transactions (id, schema_id, status, processing_time, matched_rules, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7),
		       ($8, $9, $10, $11, $12, $13, $14)
	`, transactionID1, nil, "approved", 100, matchedRulesJSON, metadataJSON, time.Now(),
		transactionID2, nil, "approved", 200, matchedRulesJSON, metadataJSON, time.Now())
	require.NoError(t, err)

	result, err := repo.ListTransactions(ctx, 10, 0)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestIntegration_TransactionsRepository_ListTransactions_WithLimit_RespectsLimit(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	matchedRulesJSON, _ := json.Marshal([]string{})
	metadataJSON, _ := json.Marshal(map[string]any{})

	for i := 0; i < 5; i++ {
		_, err := testDB.Postgres.Exec(ctx, `
			INSERT INTO transactions (id, schema_id, status, processing_time, matched_rules, metadata, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, uuid.New(), nil, "approved", 100, matchedRulesJSON, metadataJSON, time.Now())
		require.NoError(t, err)
	}

	result, err := repo.ListTransactions(ctx, 3, 0)

	require.NoError(t, err)
	assert.LessOrEqual(t, len(result), 3)
}

func TestIntegration_TransactionsRepository_ListTransactions_WithOffset_RespectsOffset(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	matchedRulesJSON, _ := json.Marshal([]string{})
	metadataJSON, _ := json.Marshal(map[string]any{})

	transactionID1 := uuid.New()
	transactionID2 := uuid.New()
	transactionID3 := uuid.New()

	now := time.Now()
	_, err := testDB.Postgres.Exec(ctx, `
		INSERT INTO transactions (id, schema_id, status, processing_time, matched_rules, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7),
		       ($8, $9, $10, $11, $12, $13, $14),
		       ($15, $16, $17, $18, $19, $20, $21)
	`, transactionID1, nil, "approved", 100, matchedRulesJSON, metadataJSON, now.Add(3*time.Second),
		transactionID2, nil, "approved", 200, matchedRulesJSON, metadataJSON, now.Add(2*time.Second),
		transactionID3, nil, "approved", 300, matchedRulesJSON, metadataJSON, now.Add(1*time.Second))
	require.NoError(t, err)

	firstPage, err := repo.ListTransactions(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, firstPage, 2)

	secondPage, err := repo.ListTransactions(ctx, 2, 2)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(secondPage), 1)
}

func TestIntegration_TransactionsRepository_ListTransactions_Empty_ReturnsEmpty(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := transactions.NewPostgresRepository(testDB.Postgres)
	ctx := context.Background()

	result, err := repo.ListTransactions(ctx, 10, 0)

	require.NoError(t, err)
	assert.NotNil(t, result)
}
