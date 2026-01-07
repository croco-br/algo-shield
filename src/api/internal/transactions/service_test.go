package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_ProcessTransaction_WhenSuccess_ThenPushesToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := models.Event{
		"external_id": "ext-123",
		"amount":      100.50,
		"currency":    "USD",
		"origin":      "account1",
		"destination": "account2",
		"type":        "transfer",
	}
	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	cmd := redis.NewIntCmd(context.Background())
	mockQueue.EXPECT().LPush(gomock.Any(), "transaction:queue", gomock.Any()).Return(cmd)
	service := NewService(mockRepo, mockQueue)

	err := service.ProcessTransaction(context.Background(), event)

	assert.NoError(t, err)
}

func Test_Service_ProcessTransaction_WhenQueueFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := models.Event{
		"external_id": "ext-123",
		"amount":      100.50,
		"currency":    "USD",
	}
	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	cmd := redis.NewIntCmd(context.Background())
	cmd.SetErr(errors.New("queue error"))
	mockQueue.EXPECT().LPush(gomock.Any(), "transaction:queue", gomock.Any()).Return(cmd)
	service := NewService(mockRepo, mockQueue)

	err := service.ProcessTransaction(context.Background(), event)

	assert.Error(t, err)
}

func Test_Service_GetTransaction_WhenExists_ThenReturnsTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	now := time.Now()
	expectedTx := &models.Transaction{
		ID:             txID,
		ExternalID:     "ext-123",
		Amount:         100.50,
		Currency:       "USD",
		Origin:         "account1",
		Destination:    "account2",
		Type:           "transfer",
		Status:         "approved",
		ProcessingTime: 10,
		CreatedAt:      now,
		ProcessedAt:    &now,
	}
	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().GetTransaction(gomock.Any(), txID).Return(expectedTx, nil)
	service := NewService(mockRepo, mockQueue)

	tx, err := service.GetTransaction(context.Background(), txID)

	require.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
}

func Test_Service_GetTransaction_WhenNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().GetTransaction(gomock.Any(), txID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo, mockQueue)

	tx, err := service.GetTransaction(context.Background(), txID)

	assert.Nil(t, tx)
	assert.Error(t, err)
}

func Test_Service_ListTransactions_WhenSuccess_ThenReturnsTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTxs := []models.Transaction{
		{
			ID:         uuid.New(),
			ExternalID: "ext-1",
			Amount:     100.50,
			Currency:   "USD",
			Status:     "approved",
		},
		{
			ID:         uuid.New(),
			ExternalID: "ext-2",
			Amount:     200.75,
			Currency:   "EUR",
			Status:     "rejected",
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return(expectedTxs, nil)
	service := NewService(mockRepo, mockQueue)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expectedTxs, txs)
}

func Test_Service_ListTransactions_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, mockQueue)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	assert.Nil(t, txs)
	assert.Error(t, err)
}

func Test_Service_ListTransactions_WhenEmptyResult_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return([]models.Transaction{}, nil)
	service := NewService(mockRepo, mockQueue)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Empty(t, txs)
}

func Test_Service_ListTransactions_WhenDifferentPagination_ThenPassesCorrectParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockQueue := NewMockQueuePusher(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 50, 100).Return([]models.Transaction{}, nil)
	service := NewService(mockRepo, mockQueue)

	_, err := service.ListTransactions(context.Background(), 50, 100)

	assert.NoError(t, err)
}
