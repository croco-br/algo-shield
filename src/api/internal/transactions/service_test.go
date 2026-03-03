package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_ProcessTransaction_WhenSuccess_ThenEnqueuesTask(t *testing.T) {
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
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)

	taskInfo := &asynq.TaskInfo{
		ID:    "task-123",
		Queue: "default",
	}

	mockEnqueuer.EXPECT().
		EnqueueTransactionWithPriority(gomock.Any(), event, "default").
		Return(taskInfo, nil)

	service := NewService(mockRepo, mockEnqueuer)

	result, err := service.ProcessTransaction(context.Background(), event)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "task-123", result.ID)
	assert.Equal(t, "default", result.Queue)
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
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)

	mockEnqueuer.EXPECT().
		EnqueueTransactionWithPriority(gomock.Any(), event, "default").
		Return(nil, errors.New("queue error"))

	service := NewService(mockRepo, mockEnqueuer)

	result, err := service.ProcessTransaction(context.Background(), event)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_Service_GetTransaction_WhenExists_ThenReturnsTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	now := time.Now()
	expectedTx := &models.Transaction{
		ID:             txID,
		Status:         models.StatusApproved,
		ProcessingTime: 10,
		Metadata: map[string]any{
			"external_id": "ext-123",
			"amount":      100.50,
			"currency":    "USD",
			"origin":      "account1",
			"destination": "account2",
			"type":        "transfer",
		},
		CreatedAt:   now,
		ProcessedAt: &now,
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetTransaction(gomock.Any(), txID).Return(expectedTx, nil)
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.GetTransaction(context.Background(), txID)

	require.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
}

func Test_Service_GetTransaction_WhenNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetTransaction(gomock.Any(), txID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.GetTransaction(context.Background(), txID)

	assert.Nil(t, tx)
	assert.Error(t, err)
}

func Test_Service_ListTransactions_WhenSuccess_ThenReturnsTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedTxs := []models.Transaction{
		{
			ID:     uuid.New(),
			Status: models.StatusApproved,
			Metadata: map[string]any{
				"external_id": "ext-1",
				"amount":      100.50,
				"currency":    "USD",
			},
		},
		{
			ID:     uuid.New(),
			Status: models.StatusRejected,
			Metadata: map[string]any{
				"external_id": "ext-2",
				"amount":      200.75,
				"currency":    "EUR",
			},
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return(expectedTxs, nil)
	service := NewService(mockRepo, mockEnqueuer)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expectedTxs, txs)
}

func Test_Service_ListTransactions_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, mockEnqueuer)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	assert.Nil(t, txs)
	assert.Error(t, err)
}

func Test_Service_ListTransactions_WhenEmptyResult_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 10, 0).Return([]models.Transaction{}, nil)
	service := NewService(mockRepo, mockEnqueuer)

	txs, err := service.ListTransactions(context.Background(), 10, 0)

	require.NoError(t, err)
	assert.Empty(t, txs)
}

func Test_Service_ListTransactions_WhenDifferentPagination_ThenPassesCorrectParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 50, 100).Return([]models.Transaction{}, nil)
	service := NewService(mockRepo, mockEnqueuer)

	_, err := service.ListTransactions(context.Background(), 50, 100)

	assert.NoError(t, err)
}

func Test_Service_ListTransactionsWithFilter_WhenSuccess_ThenReturnsResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	status := models.StatusApproved
	filter := TransactionFilter{
		Status: &status,
	}
	expectedTxs := []models.Transaction{
		{ID: uuid.New(), Status: models.StatusApproved},
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactionsWithFilter(gomock.Any(), filter, 10, 0).Return(expectedTxs, 1, nil)
	service := NewService(mockRepo, mockEnqueuer)

	txs, total, err := service.ListTransactionsWithFilter(context.Background(), filter, 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expectedTxs, txs)
	assert.Equal(t, 1, total)
}

func Test_Service_ListTransactionsWithFilter_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filter := TransactionFilter{}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ListTransactionsWithFilter(gomock.Any(), filter, 10, 0).Return(nil, 0, errors.New("database error"))
	service := NewService(mockRepo, mockEnqueuer)

	txs, total, err := service.ListTransactionsWithFilter(context.Background(), filter, 10, 0)

	assert.Nil(t, txs)
	assert.Equal(t, 0, total)
	assert.Error(t, err)
}

func Test_Service_ApproveTransaction_WhenSuccess_ThenReturnsTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	approvedTx := &models.Transaction{
		ID:     txID,
		Status: models.StatusApproved,
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ApproveTransaction(gomock.Any(), txID).Return(approvedTx, nil)
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.ApproveTransaction(context.Background(), txID)

	require.NoError(t, err)
	assert.Equal(t, approvedTx, tx)
}

func Test_Service_ApproveTransaction_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().ApproveTransaction(gomock.Any(), txID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.ApproveTransaction(context.Background(), txID)

	assert.Nil(t, tx)
	assert.Error(t, err)
}

func Test_Service_RejectTransaction_WhenSuccess_ThenReturnsTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	rejectedTx := &models.Transaction{
		ID:     txID,
		Status: models.StatusRejected,
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().RejectTransaction(gomock.Any(), txID).Return(rejectedTx, nil)
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.RejectTransaction(context.Background(), txID)

	require.NoError(t, err)
	assert.Equal(t, rejectedTx, tx)
}

func Test_Service_RejectTransaction_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	txID := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().RejectTransaction(gomock.Any(), txID).Return(nil, errors.New("not found"))
	service := NewService(mockRepo, mockEnqueuer)

	tx, err := service.RejectTransaction(context.Background(), txID)

	assert.Nil(t, tx)
	assert.Error(t, err)
}

func Test_Service_ProcessTransaction_WhenPriorityInEvent_ThenUsesPriority(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event := models.Event{
		"external_id": "ext-123",
		"amount":      100.50,
		"priority":    "critical",
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)

	taskInfo := &asynq.TaskInfo{
		ID:    "task-456",
		Queue: "critical",
	}

	mockEnqueuer.EXPECT().
		EnqueueTransactionWithPriority(gomock.Any(), event, "critical").
		Return(taskInfo, nil)

	service := NewService(mockRepo, mockEnqueuer)

	result, err := service.ProcessTransaction(context.Background(), event)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "task-456", result.ID)
}
