package transactions

import (
	"context"
	"errors"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
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
		Status:         "APPROVED",
		ProcessingTime: 50,
		MatchedRules:   []string{"rule-1"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, "tx-123", txn.ExternalID)
			assert.Equal(t, 100.50, txn.Amount)
			assert.Equal(t, "USD", txn.Currency)
			assert.Equal(t, "account-1", txn.Origin)
			assert.Equal(t, "account-2", txn.Destination)
			assert.Equal(t, "transfer", txn.Type)
			assert.Equal(t, models.TransactionStatus("APPROVED"), txn.Status)
			assert.NotNil(t, txn.ProcessedAt)
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
		Status:         "APPROVED",
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
		Status:         "BLOCKED",
		ProcessingTime: 25,
		MatchedRules:   []string{"fraud-rule-1"},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, "tx-456", txn.ExternalID)
			assert.Equal(t, 250.75, txn.Amount)
			assert.Equal(t, "EUR", txn.Currency)
			assert.Equal(t, "acc-3", txn.Origin)
			assert.Equal(t, "acc-4", txn.Destination)
			assert.Equal(t, "payment", txn.Type)
			assert.Equal(t, models.TransactionStatus("BLOCKED"), txn.Status)
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
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
		Status:         "PENDING",
		ProcessingTime: 10,
		MatchedRules:   []string{},
	}

	mockEvaluator.EXPECT().
		Evaluate(ctx, event).
		Return(expectedResult, nil)

	mockRepo.EXPECT().
		SaveTransaction(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, txn *models.Transaction) error {
			assert.Equal(t, "", txn.ExternalID)
			assert.Equal(t, 0.0, txn.Amount)
			assert.Equal(t, "", txn.Currency)
			assert.Equal(t, "", txn.Origin)
			assert.Equal(t, "", txn.Destination)
			assert.Equal(t, "", txn.Type)
			assert.NotNil(t, txn.Metadata)
			assert.Empty(t, txn.Metadata)
			return nil
		})

	err := service.ProcessTransaction(ctx, event)

	require.NoError(t, err)
}

func Test_ExtractStringFromEvent_WhenFieldExists_ThenReturnsValue(t *testing.T) {
	event := models.Event{"field1": "value1"}

	result := extractStringFromEvent(event, "field1")

	assert.Equal(t, "value1", result)
}

func Test_ExtractStringFromEvent_WhenFieldDoesNotExist_ThenReturnsEmpty(t *testing.T) {
	event := models.Event{}

	result := extractStringFromEvent(event, "field1")

	assert.Equal(t, "", result)
}

func Test_ExtractStringFromEvent_WhenFieldIsNotString_ThenReturnsEmpty(t *testing.T) {
	event := models.Event{"field1": 123}

	result := extractStringFromEvent(event, "field1")

	assert.Equal(t, "", result)
}

func Test_ExtractStringFromEvent_WithFallbacks_ThenUsesFirstMatch(t *testing.T) {
	event := models.Event{"field2": "value2"}

	result := extractStringFromEvent(event, "field1", "field2", "field3")

	assert.Equal(t, "value2", result)
}

func Test_ExtractFloat64FromEvent_WhenFieldExists_ThenReturnsValue(t *testing.T) {
	event := models.Event{"amount": 123.45}

	result := extractFloat64FromEvent(event, "amount")

	assert.Equal(t, 123.45, result)
}

func Test_ExtractFloat64FromEvent_WhenFieldIsInt_ThenConvertsToFloat(t *testing.T) {
	event := models.Event{"amount": 100}

	result := extractFloat64FromEvent(event, "amount")

	assert.Equal(t, 100.0, result)
}

func Test_ExtractFloat64FromEvent_WhenFieldDoesNotExist_ThenReturnsZero(t *testing.T) {
	event := models.Event{}

	result := extractFloat64FromEvent(event, "amount")

	assert.Equal(t, 0.0, result)
}

func Test_ToFloat64_WithFloat64_ThenReturnsValue(t *testing.T) {
	value, ok := toFloat64(123.45)

	assert.True(t, ok)
	assert.Equal(t, 123.45, value)
}

func Test_ToFloat64_WithFloat32_ThenConvertsToFloat64(t *testing.T) {
	value, ok := toFloat64(float32(123.45))

	assert.True(t, ok)
	assert.InDelta(t, 123.45, value, 0.01)
}

func Test_ToFloat64_WithInt_ThenConvertsToFloat64(t *testing.T) {
	value, ok := toFloat64(100)

	assert.True(t, ok)
	assert.Equal(t, 100.0, value)
}

func Test_ToFloat64_WithInt64_ThenConvertsToFloat64(t *testing.T) {
	value, ok := toFloat64(int64(200))

	assert.True(t, ok)
	assert.Equal(t, 200.0, value)
}

func Test_ToFloat64_WithInt32_ThenConvertsToFloat64(t *testing.T) {
	value, ok := toFloat64(int32(300))

	assert.True(t, ok)
	assert.Equal(t, 300.0, value)
}

func Test_ToFloat64_WithString_ThenReturnsFalse(t *testing.T) {
	value, ok := toFloat64("not a number")

	assert.False(t, ok)
	assert.Equal(t, 0.0, value)
}
