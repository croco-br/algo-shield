package schemas

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_CreateSchema_WhenValidRequest_ThenReturnsCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas", handler.CreateSchema)

	req := CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test description",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	body, _ := json.Marshal(req)
	expected := &EventSchema{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SampleJSON:  req.SampleJSON,
	}
	mockService.EXPECT().Create(gomock.Any(), &req).Return(expected, nil)

	httpReq := httptest.NewRequest("POST", "/schemas", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func Test_Handler_CreateSchema_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas", handler.CreateSchema)

	httpReq := httptest.NewRequest("POST", "/schemas", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_CreateSchema_WhenEmptySampleJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas", handler.CreateSchema)

	req := CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test description",
		SampleJSON:  map[string]any{},
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest("POST", "/schemas", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_CreateSchema_WhenNameExists_ThenReturnsConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas", handler.CreateSchema)

	req := CreateSchemaRequest{
		Name:        "existing-schema",
		Description: "Test description",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	body, _ := json.Marshal(req)
	mockService.EXPECT().Create(gomock.Any(), &req).Return(nil, ErrSchemaNameExists)

	httpReq := httptest.NewRequest("POST", "/schemas", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_Handler_CreateSchema_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas", handler.CreateSchema)

	req := CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test description",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	body, _ := json.Marshal(req)
	mockService.EXPECT().Create(gomock.Any(), &req).Return(nil, errors.New("database error"))

	httpReq := httptest.NewRequest("POST", "/schemas", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_GetSchema_WhenValidID_ThenReturnsSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas/:id", handler.GetSchema)

	id := uuid.New()
	expected := &EventSchema{
		ID:   id,
		Name: "test-schema",
	}
	mockService.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

	httpReq := httptest.NewRequest("GET", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_GetSchema_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas/:id", handler.GetSchema)

	httpReq := httptest.NewRequest("GET", "/schemas/invalid-uuid", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_GetSchema_WhenNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas/:id", handler.GetSchema)

	id := uuid.New()
	mockService.EXPECT().GetByID(gomock.Any(), id).Return(nil, ErrSchemaNotFound)

	httpReq := httptest.NewRequest("GET", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_GetSchema_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas/:id", handler.GetSchema)

	id := uuid.New()
	mockService.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("database error"))

	httpReq := httptest.NewRequest("GET", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_ListSchemas_WhenSuccess_ThenReturnsSchemas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas", handler.ListSchemas)

	expected := []EventSchema{
		{ID: uuid.New(), Name: "schema1"},
		{ID: uuid.New(), Name: "schema2"},
	}
	mockService.EXPECT().List(gomock.Any()).Return(expected, nil)

	httpReq := httptest.NewRequest("GET", "/schemas", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result SchemaListResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.Len(t, result.Schemas, 2)
}

func Test_Handler_ListSchemas_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/schemas", handler.ListSchemas)

	mockService.EXPECT().List(gomock.Any()).Return(nil, errors.New("database error"))

	httpReq := httptest.NewRequest("GET", "/schemas", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenValidRequest_ThenReturnsUpdatedSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	id := uuid.New()
	req := UpdateSchemaRequest{
		Name:        "updated-name",
		Description: "Updated description",
	}
	body, _ := json.Marshal(req)
	expected := &EventSchema{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	mockService.EXPECT().Update(gomock.Any(), id, &req).Return(expected, nil)

	httpReq := httptest.NewRequest("PUT", "/schemas/"+id.String(), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	req := UpdateSchemaRequest{Name: "updated-name"}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest("PUT", "/schemas/invalid-uuid", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	id := uuid.New()

	httpReq := httptest.NewRequest("PUT", "/schemas/"+id.String(), bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	id := uuid.New()
	req := UpdateSchemaRequest{Name: "updated-name"}
	body, _ := json.Marshal(req)
	mockService.EXPECT().Update(gomock.Any(), id, &req).Return(nil, ErrSchemaNotFound)

	httpReq := httptest.NewRequest("PUT", "/schemas/"+id.String(), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenNameExists_ThenReturnsConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	id := uuid.New()
	req := UpdateSchemaRequest{Name: "existing-name"}
	body, _ := json.Marshal(req)
	mockService.EXPECT().Update(gomock.Any(), id, &req).Return(nil, ErrSchemaNameExists)

	httpReq := httptest.NewRequest("PUT", "/schemas/"+id.String(), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_Handler_UpdateSchema_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/schemas/:id", handler.UpdateSchema)

	id := uuid.New()
	req := UpdateSchemaRequest{Name: "updated-name"}
	body, _ := json.Marshal(req)
	mockService.EXPECT().Update(gomock.Any(), id, &req).Return(nil, errors.New("database error"))

	httpReq := httptest.NewRequest("PUT", "/schemas/"+id.String(), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_DeleteSchema_WhenValidID_ThenReturnsNoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Delete("/schemas/:id", handler.DeleteSchema)

	id := uuid.New()
	mockService.EXPECT().Delete(gomock.Any(), id).Return(nil)

	httpReq := httptest.NewRequest("DELETE", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func Test_Handler_DeleteSchema_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Delete("/schemas/:id", handler.DeleteSchema)

	httpReq := httptest.NewRequest("DELETE", "/schemas/invalid-uuid", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_DeleteSchema_WhenNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Delete("/schemas/:id", handler.DeleteSchema)

	id := uuid.New()
	mockService.EXPECT().Delete(gomock.Any(), id).Return(ErrSchemaNotFound)

	httpReq := httptest.NewRequest("DELETE", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_DeleteSchema_WhenHasRules_ThenReturnsConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Delete("/schemas/:id", handler.DeleteSchema)

	id := uuid.New()
	mockService.EXPECT().Delete(gomock.Any(), id).Return(ErrSchemaHasRules)
	mockService.EXPECT().GetRulesReferencingSchema(gomock.Any(), id).Return([]string{"rule1", "rule2"}, nil)

	httpReq := httptest.NewRequest("DELETE", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func Test_Handler_DeleteSchema_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Delete("/schemas/:id", handler.DeleteSchema)

	id := uuid.New()
	mockService.EXPECT().Delete(gomock.Any(), id).Return(errors.New("database error"))

	httpReq := httptest.NewRequest("DELETE", "/schemas/"+id.String(), nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_ParseSchema_WhenValidID_ThenReturnsUpdatedSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas/:id/parse", handler.ParseSchema)

	id := uuid.New()
	expected := &EventSchema{
		ID:   id,
		Name: "test-schema",
	}
	mockService.EXPECT().ParseSampleJSON(gomock.Any(), id).Return(expected, nil)

	httpReq := httptest.NewRequest("POST", "/schemas/"+id.String()+"/parse", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_ParseSchema_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas/:id/parse", handler.ParseSchema)

	httpReq := httptest.NewRequest("POST", "/schemas/invalid-uuid/parse", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_ParseSchema_WhenNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas/:id/parse", handler.ParseSchema)

	id := uuid.New()
	mockService.EXPECT().ParseSampleJSON(gomock.Any(), id).Return(nil, ErrSchemaNotFound)

	httpReq := httptest.NewRequest("POST", "/schemas/"+id.String()+"/parse", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_ParseSchema_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockServiceInterface(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Post("/schemas/:id/parse", handler.ParseSchema)

	id := uuid.New()
	mockService.EXPECT().ParseSampleJSON(gomock.Any(), id).Return(nil, errors.New("database error"))

	httpReq := httptest.NewRequest("POST", "/schemas/"+id.String()+"/parse", nil)
	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
