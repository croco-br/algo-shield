package permissions

import (
	"bytes"
	"encoding/json"
	"errors"
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
	service := &Service{}

	handler := NewHandler(service)

	require.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

func Test_Handler_ListUsers_WhenSuccess_ThenReturnsUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	expectedUsers := []models.User{
		{ID: uuid.New(), Email: "user1@example.com", Name: "User 1"},
		{ID: uuid.New(), Email: "user2@example.com", Name: "User 2"},
	}
	mockService.EXPECT().
		ListUsers(gomock.Any()).
		Return(expectedUsers, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/users", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_ListUsers_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	mockService.EXPECT().
		ListUsers(gomock.Any()).
		Return(nil, errors.New("database error"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/users", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func Test_Handler_GetUser_WhenValidID_ThenReturnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	userID := uuid.New()
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(expectedUser, nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_GetUser_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/users/invalid-uuid", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_GetUser_WhenUserNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	userID := uuid.New()
	mockService.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(nil, errors.New("not found"))
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Get("/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_UpdateUserActive_WhenValidRequest_ThenUpdatesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	targetUserID := uuid.New()
	mockService.EXPECT().
		UpdateUserActive(gomock.Any(), gomock.Any(), targetUserID, false).
		Return(nil)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/users/:id/active", func(c *fiber.Ctx) error {
		currentUser := &models.User{
			ID:    uuid.New(),
			Email: "admin@example.com",
		}
		c.Locals("user", currentUser)
		return handler.UpdateUserActive(c)
	})
	reqBody := UpdateUserActiveRequest{Active: false}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/users/"+targetUserID.String()+"/active", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_UpdateUserActive_WhenUserNotInContext_ThenReturnsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	handler := NewHandler(mockService)
	app := fiber.New()
	app.Put("/users/:id/active", handler.UpdateUserActive)
	targetUserID := uuid.New()
	reqBody := UpdateUserActiveRequest{Active: false}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/users/"+targetUserID.String()+"/active", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_Handler_UpdateUserActive_WhenInvalidUserID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Put("/users/:id/active", func(c *fiber.Ctx) error {
		currentUser := &models.User{ID: uuid.New()}
		c.Locals("user", currentUser)
		return handler.UpdateUserActive(c)
	})

	reqBody := UpdateUserActiveRequest{Active: false}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/users/invalid-uuid/active", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_UpdateUserActive_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockPermissionsService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Put("/users/:id/active", func(c *fiber.Ctx) error {
		currentUser := &models.User{ID: uuid.New()}
		c.Locals("user", currentUser)
		return handler.UpdateUserActive(c)
	})

	targetUserID := uuid.New()
	req := httptest.NewRequest("PUT", "/users/"+targetUserID.String()+"/active", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
