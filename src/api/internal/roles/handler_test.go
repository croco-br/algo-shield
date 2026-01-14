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

func Test_Handler_GetRole(t *testing.T) {
	testCases := []struct {
		name           string
		roleID         string
		setupMock      func(*MockService, uuid.UUID)
		expectedStatus int
		validateBody   func(*testing.T, []byte)
	}{
		{
			name:   "valid ID returns role",
			roleID: uuid.New().String(),
			setupMock: func(m *MockService, id uuid.UUID) {
				expectedRole := &models.Role{
					ID:          id,
					Name:        "admin",
					Description: "Administrator",
				}
				m.EXPECT().GetRoleByID(gomock.Any(), id).Return(expectedRole, nil)
			},
			expectedStatus: fiber.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var result models.Role
				err := json.Unmarshal(body, &result)
				require.NoError(t, err)
				assert.Equal(t, "admin", result.Name)
			},
		},
		{
			name:           "invalid ID returns bad request",
			roleID:         "invalid-uuid",
			setupMock:      func(m *MockService, id uuid.UUID) {},
			expectedStatus: fiber.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var result map[string]string
				err := json.Unmarshal(body, &result)
				require.NoError(t, err)
				assert.Equal(t, "Invalid role ID", result["error"])
			},
		},
		{
			name:   "role not found returns not found",
			roleID: uuid.New().String(),
			setupMock: func(m *MockService, id uuid.UUID) {
				m.EXPECT().GetRoleByID(gomock.Any(), id).Return(nil, errors.New("not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			validateBody: func(t *testing.T, body []byte) {
				var result map[string]string
				err := json.Unmarshal(body, &result)
				require.NoError(t, err)
				assert.Equal(t, "Role not found", result["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockService(ctrl)
			handler := NewHandler(mockService)
			app := fiber.New()
			app.Get("/roles/:id", handler.GetRole)

			roleID, _ := uuid.Parse(tc.roleID)
			tc.setupMock(mockService, roleID)

			req := httptest.NewRequest("GET", "/roles/"+tc.roleID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			tc.validateBody(t, body)
		})
	}
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

func Test_Handler_RemoveRole(t *testing.T) {
	testCases := []struct {
		name           string
		userID         string
		roleID         string
		setupMock      func(*MockService, uuid.UUID, uuid.UUID)
		expectedStatus int
	}{
		{
			name:   "valid request removes role successfully",
			userID: uuid.New().String(),
			roleID: uuid.New().String(),
			setupMock: func(m *MockService, userID, roleID uuid.UUID) {
				m.EXPECT().RemoveRole(gomock.Any(), userID, roleID).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "invalid user ID returns bad request",
			userID:         "invalid-uuid",
			roleID:         uuid.New().String(),
			setupMock:      func(m *MockService, userID, roleID uuid.UUID) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid role ID returns bad request",
			userID:         uuid.New().String(),
			roleID:         "invalid-uuid",
			setupMock:      func(m *MockService, userID, roleID uuid.UUID) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:   "service failure returns internal error",
			userID: uuid.New().String(),
			roleID: uuid.New().String(),
			setupMock: func(m *MockService, userID, roleID uuid.UUID) {
				m.EXPECT().RemoveRole(gomock.Any(), userID, roleID).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockService(ctrl)
			handler := NewHandler(mockService)
			app := fiber.New()
			app.Delete("/users/:userId/roles/:roleId", handler.RemoveRole)

			userID, _ := uuid.Parse(tc.userID)
			roleID, _ := uuid.Parse(tc.roleID)
			tc.setupMock(mockService, userID, roleID)

			req := httptest.NewRequest("DELETE", "/users/"+tc.userID+"/roles/"+tc.roleID, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
