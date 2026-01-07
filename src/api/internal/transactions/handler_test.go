package transactions

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	mockService.EXPECT().
		ProcessTransaction(gomock.Any(), gomock.Any()).
		Return(nil)

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
		Return(errors.New("queue error"))

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
		ID:         transactionID,
		ExternalID: "tx-123",
		Amount:     100.0,
		Currency:   "USD",
		Status:     "APPROVED",
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

	assert.Equal(t, expectedTransaction.ExternalID, result.ExternalID)
	assert.Equal(t, expectedTransaction.Amount, result.Amount)
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
		{ID: uuid.New(), ExternalID: "tx-1", Amount: 100.0},
		{ID: uuid.New(), ExternalID: "tx-2", Amount: 200.0},
	}

	mockService.EXPECT().
		ListTransactions(gomock.Any(), 50, 0).
		Return(expectedTransactions, nil)

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
}

func Test_Handler_ListTransactions_WhenCustomPagination_ThenUsesParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockTransactionService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/transactions", handler.ListTransactions)

	expectedTransactions := []models.Transaction{
		{ID: uuid.New(), ExternalID: "tx-1", Amount: 100.0},
	}

	mockService.EXPECT().
		ListTransactions(gomock.Any(), 10, 5).
		Return(expectedTransactions, nil)

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
		ListTransactions(gomock.Any(), 50, 0).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/transactions", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
