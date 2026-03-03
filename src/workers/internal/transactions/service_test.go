package transactions

import (
	"context"
	"errors"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_ProcessTransaction_WhenValidEvent_ThenSavesTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{
		"external_id": "tx-123",
		"amount":      100.50,
		"currency":    "USD",
		"origin":      "account-1",
		"destination": "account-2",
		"type":        "transfer",
		"metadata": map[string]any{
			"ip_address": "192.168.1.1",
		},
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusApproved,
		ProcessingTime: 50,
		MatchedRules:   []string{"rule-1"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, models.StatusApproved, txn.Status)
			assert.NotNil(t, txn.ProcessedAt)
			assert.NotNil(t, txn.Metadata)
			assert.Equal(t, "tx-123", txn.Metadata["external_id"])
			assert.Equal(t, 100.50, txn.Metadata["amount"])
			assert.Equal(t, "USD", txn.Metadata["currency"])
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenEvaluatorFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{"external_id": "tx-123"}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(nil, errors.New("evaluation error"))

	err := service.ProcessTransaction(ctx, event)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "evaluation error")
}

func Test_Service_ProcessTransaction_WhenSaveFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{"external_id": "tx-123"}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusApproved,
		ProcessingTime: 50,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		Return(errors.New("database error"))

	err := service.ProcessTransaction(ctx, event)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func Test_Service_ProcessTransaction_WhenAlternativeFieldNames_ThenExtractsCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{
		"id":            "tx-456",
		"value":         250.75,
		"currency_code": "EUR",
		"from_account":  "acc-3",
		"recipient_id":  "acc-4",
		"event_type":    "payment",
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusRejected,
		ProcessingTime: 25,
		MatchedRules:   []string{"fraud-rule-1"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, models.StatusRejected, txn.Status)
			assert.NotNil(t, txn.Metadata)
			assert.Equal(t, "tx-456", txn.Metadata["id"])
			assert.Equal(t, 250.75, txn.Metadata["value"])
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenSyntheticEvent_ThenSavesSyntheticTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{
		"amount":     500.0,
		"currency":   "BRL",
		"_synthetic": true,
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusRejected,
		ProcessingTime: 30,
		MatchedRules:   []string{"high-amount"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveSyntheticTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			// Synthetic transactions always get pending status
			assert.Equal(t, models.StatusPending, txn.Status)
			assert.Nil(t, txn.ProcessedAt)
			assert.Equal(t, int64(0), txn.ProcessingTime)
			assert.Empty(t, txn.MatchedRules)
			// _synthetic should be stripped from metadata
			assert.NotContains(t, txn.Metadata, "_synthetic")
			assert.Equal(t, 500.0, txn.Metadata["amount"])
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenSchemaIDPresent_ThenExtractsAndRemovesFromMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	schemaID := uuid.New()
	event := models.Event{
		"amount":     100.0,
		"_schema_id": schemaID.String(),
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusApproved,
		ProcessingTime: 10,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			// SchemaID should be set
			require.NotNil(t, txn.SchemaID)
			assert.Equal(t, schemaID, *txn.SchemaID)
			// _schema_id should be stripped from metadata
			assert.NotContains(t, txn.Metadata, "_schema_id")
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenInvalidSchemaID_ThenSchemaIDIsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{
		"amount":     100.0,
		"_schema_id": "not-a-valid-uuid",
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusApproved,
		ProcessingTime: 10,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			// SchemaID should remain nil for invalid UUID
			assert.Nil(t, txn.SchemaID)
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenSyntheticWithSchemaID_ThenSavesSyntheticWithSchemaID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	schemaID := uuid.New()
	event := models.Event{
		"amount":     200.0,
		"_synthetic": true,
		"_schema_id": schemaID.String(),
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusRejected,
		ProcessingTime: 20,
		MatchedRules:   []string{"rule-1"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveSyntheticTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			// Should have schema ID
			require.NotNil(t, txn.SchemaID)
			assert.Equal(t, schemaID, *txn.SchemaID)
			// Should be pending (synthetic)
			assert.Equal(t, models.StatusPending, txn.Status)
			// Both flags stripped from metadata
			assert.NotContains(t, txn.Metadata, "_synthetic")
			assert.NotContains(t, txn.Metadata, "_schema_id")
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenSyntheticSaveFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{
		"amount":     100.0,
		"_synthetic": true,
	}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusApproved,
		ProcessingTime: 10,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveSyntheticTransaction(ctx, gomock.Any()).
		Return(errors.New("synthetic save failed"))

	err := service.ProcessTransaction(ctx, event)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "synthetic save failed")
}

func Test_Service_ProcessTransaction_WhenMissingFields_ThenUsesDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEvaluator := NewMockRuleEvaluator(ctrl)

	service := NewService(mockRepo, mockEvaluator)

	ctx := context.Background()
	event := models.Event{}

	expectedResult := &models.TransactionResult{
		Status:         models.StatusPending,
		ProcessingTime: 10,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, models.StatusPending, txn.Status)
			assert.NotNil(t, txn.Metadata)
			assert.Empty(t, txn.Metadata)
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}
