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

	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO transactions (id, external_id, amount, currency, origin, destination, type, status, processing_time, matched_rules, metadata, created_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, transactionID, externalID, 100.50, "USD", "account1", "account2", "transfer", "approved", 150, matchedRulesJSON, metadataJSON, time.Now(), time.Now())
	require.NoError(t, err)

	result, err := repo.GetTransaction(ctx, transactionID)

	require.NoError(t, err)
	assert.Equal(t, transactionID, result.ID)
	assert.Equal(t, externalID, result.ExternalID)
	assert.Equal(t, 100.50, result.Amount)
	assert.Equal(t, "USD", result.Currency)
	assert.Equal(t, "account1", result.Origin)
	assert.Equal(t, "account2", result.Destination)
	assert.Equal(t, "transfer", result.Type)
	assert.Equal(t, models.StatusApproved, result.Status)
	assert.Equal(t, int64(150), result.ProcessingTime)
	assert.Len(t, result.MatchedRules, 2)
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
		INSERT INTO transactions (id, external_id, amount, currency, origin, destination, type, status, processing_time, matched_rules, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12),
		       ($13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`, transactionID1, "ext-1", 100.0, "USD", "acc1", "acc2", "transfer", "approved", 100, matchedRulesJSON, metadataJSON, time.Now(),
		transactionID2, "ext-2", 200.0, "EUR", "acc3", "acc4", "payment", "approved", 200, matchedRulesJSON, metadataJSON, time.Now())
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
			INSERT INTO transactions (id, external_id, amount, currency, origin, destination, type, status, processing_time, matched_rules, metadata, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, uuid.New(), "ext-"+strconv.Itoa(i), 100.0, "USD", "acc1", "acc2", "transfer", "approved", 100, matchedRulesJSON, metadataJSON, time.Now())
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
		INSERT INTO transactions (id, external_id, amount, currency, origin, destination, type, status, processing_time, matched_rules, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12),
		       ($13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24),
		       ($25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36)
	`, transactionID1, "ext-1", 100.0, "USD", "acc1", "acc2", "transfer", "approved", 100, matchedRulesJSON, metadataJSON, now.Add(3*time.Second),
		transactionID2, "ext-2", 200.0, "EUR", "acc3", "acc4", "payment", "approved", 200, matchedRulesJSON, metadataJSON, now.Add(2*time.Second),
		transactionID3, "ext-3", 300.0, "GBP", "acc5", "acc6", "transfer", "approved", 300, matchedRulesJSON, metadataJSON, now.Add(1*time.Second))
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
