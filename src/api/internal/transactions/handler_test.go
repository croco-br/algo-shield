package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_NewHandler_WhenCalled_ThenReturnsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)

	handler := NewHandler(mockService)

	require.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func Test_Handler_ProcessTransaction_WhenValidEvent_ThenQueuesTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/transactions", handler.ProcessTransaction)

	event := models.Event{
		"external_id": "tx-123",
		"amount":      100.0,
		"currency":    "USD",
	}

	taskInfo := &asynq.TaskInfo{
		ID:    "task-123",
		Queue: "default",
	}

	mockService.EXPECT().
		ProcessTransaction(gomock.Any(), gomock.Any()).
		Return(taskInfo, nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusAccepted, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "queued", result["status"])
	assert.Equal(t, "tx-123", result["external_id"])
	assert.Equal(t, "task-123", result["job_id"])
	assert.Equal(t, "default", result["queue"])
}

func Test_Handler_ProcessTransaction_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/transactions", handler.ProcessTransaction)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_ProcessTransaction_WhenEmptyEvent_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/transactions", handler.ProcessTransaction)

	event := models.Event{}
	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "Event must be a non-empty JSON object", result["error"])
}

func Test_Handler_ProcessTransaction_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/transactions", handler.ProcessTransaction)

	event := models.Event{
		"external_id": "tx-123",
		"amount":      100.0,
	}

	mockService.EXPECT().
		ProcessTransaction(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("queue error"))

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_GetTransaction_WhenValidID_ThenReturnsTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions/:id", handler.GetTransaction)

	transactionID := uuid.New()
	expectedTransaction := &models.Transaction{
		ID:     transactionID,
		Status: models.StatusApproved,
		Metadata: map[string]any{
			"external_id": "tx-123",
			"amount":      100.0,
			"currency":    "USD",
		},
	}

	mockService.EXPECT().
		GetTransaction(gomock.Any(), transactionID).
		Return(expectedTransaction, nil)

	req := httptest.NewRequest("GET", "/transactions/"+transactionID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result models.Transaction
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, expectedTransaction.Status, result.Status)
	assert.Equal(t, "tx-123", result.Metadata["external_id"])
	assert.Equal(t, 100.0, result.Metadata["amount"])
}

func Test_Handler_GetTransaction_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions/:id", handler.GetTransaction)

	req := httptest.NewRequest("GET", "/transactions/invalid-uuid", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Invalid transaction ID", result["error"])
}

func Test_Handler_GetTransaction_WhenNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions/:id", handler.GetTransaction)

	transactionID := uuid.New()

	mockService.EXPECT().
		GetTransaction(gomock.Any(), transactionID).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/transactions/"+transactionID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_ListTransactions_WhenSuccess_ThenReturnsTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	expectedTransactions := []models.Transaction{
		{ID: uuid.New(), Status: models.StatusApproved, Metadata: map[string]any{"external_id": "tx-1", "amount": 100.0}},
		{ID: uuid.New(), Status: models.StatusApproved, Metadata: map[string]any{"external_id": "tx-2", "amount": 200.0}},
	}

	mockService.EXPECT().
		ListTransactionsWithFilter(gomock.Any(), TransactionFilter{}, 50, 0).
		Return(expectedTransactions, 2, nil)

	req := httptest.NewRequest("GET", "/transactions", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	transactions := result["transactions"].([]interface{})
	assert.Len(t, transactions, 2)
	assert.Equal(t, float64(2), result["total"])
}

func Test_Handler_ListTransactions_WhenCustomPagination_ThenUsesParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	expectedTransactions := []models.Transaction{
		{ID: uuid.New(), Status: models.StatusApproved, Metadata: map[string]any{"external_id": "tx-1", "amount": 100.0}},
	}

	mockService.EXPECT().
		ListTransactionsWithFilter(gomock.Any(), TransactionFilter{}, 10, 5).
		Return(expectedTransactions, 1, nil)

	req := httptest.NewRequest("GET", "/transactions?limit=10&offset=5", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_ListTransactions_WhenInvalidLimit_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	req := httptest.NewRequest("GET", "/transactions?limit=0", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_ListTransactions_WhenInvalidOffset_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	req := httptest.NewRequest("GET", "/transactions?offset=-1", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_ListTransactions_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	mockService.EXPECT().
		ListTransactionsWithFilter(gomock.Any(), TransactionFilter{}, 50, 0).
		Return(nil, 0, errors.New("database error"))

	req := httptest.NewRequest("GET", "/transactions", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_ListTransactions_WhenInvalidStartDate_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	req := httptest.NewRequest("GET", "/transactions?start_date=invalid-date", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Contains(t, result["error"], "Invalid start_date format")
}

func Test_Handler_ListTransactions_WhenInvalidEndDate_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	req := httptest.NewRequest("GET", "/transactions?end_date=invalid-date", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Contains(t, result["error"], "Invalid end_date format")
}

func Test_Handler_ListTransactions_WhenValidDateFilters_ThenUsesFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 999999999, time.UTC)

	expectedTransactions := []models.Transaction{
		{ID: uuid.New(), Status: models.StatusApproved, Metadata: map[string]any{"external_id": "tx-1", "amount": 100.0}},
	}

	mockService.EXPECT().
		ListTransactionsWithFilter(gomock.Any(), gomock.Any(), 50, 0).
		DoAndReturn(func(ctx interface{}, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error) {
			assert.NotNil(t, filter.StartDate)
			assert.NotNil(t, filter.EndDate)
			assert.WithinDuration(t, startDate, *filter.StartDate, time.Second)
			assert.WithinDuration(t, endDate, *filter.EndDate, time.Second)
			return expectedTransactions, 1, nil
		})

	req := httptest.NewRequest("GET", "/transactions?start_date=2025-01-01T00:00:00Z&end_date=2025-01-31T23:59:59.999Z", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
