package roles

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

	mockService := NewMockService(ctrl)

	handler := NewHandler(mockService)

	require.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func Test_Handler_ListRoles_WhenSuccess_ThenReturnsRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "admin", Description: "Administrator"},
		{ID: uuid.New(), Name: "viewer", Description: "Viewer"},
	}
	mockService.EXPECT().
		ListRoles(gomock.Any()).
		Return(expectedRoles, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/roles", handler.ListRoles)

	req := httptest.NewRequest("GET", "/roles", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result map[string][]models.Role
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.Len(t, result["roles"], 2)
	assert.Equal(t, expectedRoles[0].Name, result["roles"][0].Name)
}

func Test_Handler_ListRoles_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/roles", handler.ListRoles)

	mockService.EXPECT().
		ListRoles(gomock.Any()).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/roles", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Failed to fetch roles", result["error"])
}

func Test_Handler_GetRole_WhenValidID_ThenReturnsRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/roles/:id", handler.GetRole)

	roleID := uuid.New()
	expectedRole := &models.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator",
	}

	mockService.EXPECT().
		GetRoleByID(gomock.Any(), roleID).
		Return(expectedRole, nil)

	req := httptest.NewRequest("GET", "/roles/"+roleID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result models.Role
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, expectedRole.Name, result.Name)
	assert.Equal(t, expectedRole.Description, result.Description)
}

func Test_Handler_GetRole_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/roles/:id", handler.GetRole)

	req := httptest.NewRequest("GET", "/roles/invalid-uuid", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Invalid role ID", result["error"])
}

func Test_Handler_GetRole_WhenRoleNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/roles/:id", handler.GetRole)

	roleID := uuid.New()

	mockService.EXPECT().
		GetRoleByID(gomock.Any(), roleID).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/roles/"+roleID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Role not found", result["error"])
}

func Test_Handler_AssignRole_WhenValidRequest_ThenAssignsRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/users/:userId/roles", handler.AssignRole)

	userID := uuid.New()
	roleID := uuid.New()

	mockService.EXPECT().
		AssignRole(gomock.Any(), userID, roleID).
		Return(nil)

	reqBody := AssignRoleRequest{RoleID: roleID}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users/"+userID.String()+"/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.Equal(t, "Role assigned successfully", result["message"])
}

func Test_Handler_AssignRole_WhenInvalidUserID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/users/:userId/roles", handler.AssignRole)

	roleID := uuid.New()
	reqBody := AssignRoleRequest{RoleID: roleID}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users/invalid-uuid/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_AssignRole_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/users/:userId/roles", handler.AssignRole)

	userID := uuid.New()
	req := httptest.NewRequest("POST", "/users/"+userID.String()+"/roles", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_AssignRole_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/users/:userId/roles", handler.AssignRole)

	userID := uuid.New()
	roleID := uuid.New()

	mockService.EXPECT().
		AssignRole(gomock.Any(), userID, roleID).
		Return(errors.New("database error"))

	reqBody := AssignRoleRequest{RoleID: roleID}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users/"+userID.String()+"/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_RemoveRole_WhenValidRequest_ThenRemovesRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Delete("/users/:userId/roles/:roleId", handler.RemoveRole)

	userID := uuid.New()
	roleID := uuid.New()

	mockService.EXPECT().
		RemoveRole(gomock.Any(), userID, roleID).
		Return(nil)

	req := httptest.NewRequest("DELETE", "/users/"+userID.String()+"/roles/"+roleID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Role removed successfully", result["message"])
}

func Test_Handler_RemoveRole_WhenInvalidUserID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Delete("/users/:userId/roles/:roleId", handler.RemoveRole)

	roleID := uuid.New()
	req := httptest.NewRequest("DELETE", "/users/invalid-uuid/roles/"+roleID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_RemoveRole_WhenInvalidRoleID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Delete("/users/:userId/roles/:roleId", handler.RemoveRole)

	userID := uuid.New()
	req := httptest.NewRequest("DELETE", "/users/"+userID.String()+"/roles/invalid-uuid", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_RemoveRole_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Delete("/users/:userId/roles/:roleId", handler.RemoveRole)

	userID := uuid.New()
	roleID := uuid.New()

	mockService.EXPECT().
		RemoveRole(gomock.Any(), userID, roleID).
		Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/users/"+userID.String()+"/roles/"+roleID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}
