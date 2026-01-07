package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
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

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)

	handler := NewHandler(mockService, mockUserService)

	require.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
	assert.Equal(t, mockUserService, handler.userService)
}

func Test_Handler_Register_WhenValidRequest_ThenReturnsUserAndToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	expectedToken := "test-token-123"
	mockService.EXPECT().
		RegisterUser(gomock.Any(), "test@example.com", "Test User", "password123").
		Return(expectedUser, expectedToken, nil)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/register", handler.Register)
	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_Register_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_Register_WhenValidationFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/register", handler.Register)
	reqBody := RegisterRequest{
		Email:    "invalid-email",
		Name:     "Test User",
		Password: "short",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_Login_WhenValidCredentials_ThenReturnsUserAndToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	expectedToken := "test-token-123"
	mockService.EXPECT().
		LoginUser(gomock.Any(), "test@example.com", "password123").
		Return(expectedUser, expectedToken, nil)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/login", handler.Login)
	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_Login_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_Login_WhenValidationFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/login", handler.Login)
	reqBody := LoginRequest{
		Email:    "invalid-email",
		Password: "short",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_GetCurrentUser_WhenUserInContext_ThenReturnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(expectedUser, nil)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Get("/auth/me", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Name:  "Test User",
		}
		c.Locals("user", user)
		return handler.GetCurrentUser(c)
	})

	req := httptest.NewRequest("GET", "/auth/me", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_GetCurrentUser_WhenUserNotInContext_ThenReturnsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Get("/auth/me", handler.GetCurrentUser)

	req := httptest.NewRequest("GET", "/auth/me", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_Handler_GetCurrentUser_WhenServiceFails_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	mockUserService.EXPECT().
		GetUserByID(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("user not found"))
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Get("/auth/me", func(c *fiber.Ctx) error {
		user := &models.User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}
		c.Locals("user", user)
		return handler.GetCurrentUser(c)
	})

	req := httptest.NewRequest("GET", "/auth/me", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Handler_Logout_WhenValidToken_ThenReturnsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	mockService.EXPECT().
		LogoutUser(gomock.Any(), "valid-token-123").
		Return(nil)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid-token-123")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	_ = json.Unmarshal(body, &result)
	assert.Equal(t, "Logged out successfully", result["message"])
}

func Test_Handler_Logout_WhenNoAuthHeader_ThenReturnsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_Handler_Logout_WhenInvalidAuthHeaderFormat_ThenReturnsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "InvalidFormat")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_Handler_Logout_WhenMissingBearerPrefix_ThenReturnsUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	handler := NewHandler(mockService, mockUserService)
	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Authorization", "Token valid-token-123")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_Handler_ValidateToken_WhenValidToken_ThenReturnsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockService.EXPECT().
		ValidateToken("valid-token").
		Return(expectedUser, nil)
	handler := NewHandler(mockService, mockUserService)

	user, err := handler.ValidateToken("valid-token")

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func Test_Handler_ValidateToken_WhenInvalidToken_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockAuthService(ctrl)
	mockUserService := NewMockUserService(ctrl)
	mockService.EXPECT().
		ValidateToken("invalid-token").
		Return(nil, apierrors.TokenInvalid())
	handler := NewHandler(mockService, mockUserService)

	user, err := handler.ValidateToken("invalid-token")

	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, apierrors.ErrTokenInvalid, err.(*apierrors.APIError).Code)
}
