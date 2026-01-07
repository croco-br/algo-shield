package groups

import (
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

func Test_Handler_ListGroups_WhenSuccess_ThenReturnsGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups", handler.ListGroups)

	expectedGroups := []models.Group{
		{ID: uuid.New(), Name: "Group 1", Description: "Description 1"},
		{ID: uuid.New(), Name: "Group 2", Description: "Description 2"},
	}

	mockService.EXPECT().
		ListGroups(gomock.Any()).
		Return(expectedGroups, nil)

	req := httptest.NewRequest("GET", "/groups", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string][]models.Group
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Len(t, result["groups"], 2)
	assert.Equal(t, expectedGroups[0].Name, result["groups"][0].Name)
}

func Test_Handler_ListGroups_WhenServiceFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups", handler.ListGroups)

	mockService.EXPECT().
		ListGroups(gomock.Any()).
		Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/groups", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Failed to fetch groups", result["error"])
}

func Test_Handler_GetGroup_WhenValidID_ThenReturnsGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups/:id", handler.GetGroup)

	groupID := uuid.New()
	expectedGroup := &models.Group{
		ID:          groupID,
		Name:        "Test Group",
		Description: "Test Description",
	}
	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "Role 1"},
		{ID: uuid.New(), Name: "Role 2"},
	}

	mockService.EXPECT().
		GetGroupByID(gomock.Any(), groupID).
		Return(expectedGroup, nil)

	mockService.EXPECT().
		LoadGroupRoles(gomock.Any(), groupID).
		Return(expectedRoles, nil)

	req := httptest.NewRequest("GET", "/groups/"+groupID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result models.Group
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, expectedGroup.Name, result.Name)
	assert.Len(t, result.Roles, 2)
}

func Test_Handler_GetGroup_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups/:id", handler.GetGroup)

	req := httptest.NewRequest("GET", "/groups/invalid-uuid", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Invalid group ID", result["error"])
}

func Test_Handler_GetGroup_WhenGroupNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups/:id", handler.GetGroup)

	groupID := uuid.New()

	mockService.EXPECT().
		GetGroupByID(gomock.Any(), groupID).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/groups/"+groupID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, "Group not found", result["error"])
}

func Test_Handler_GetGroup_WhenLoadRolesFails_ThenReturnsGroupWithoutRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockService(ctrl)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Get("/groups/:id", handler.GetGroup)

	groupID := uuid.New()
	expectedGroup := &models.Group{
		ID:          groupID,
		Name:        "Test Group",
		Description: "Test Description",
	}

	mockService.EXPECT().
		GetGroupByID(gomock.Any(), groupID).
		Return(expectedGroup, nil)

	mockService.EXPECT().
		LoadGroupRoles(gomock.Any(), groupID).
		Return(nil, errors.New("roles load failed"))

	req := httptest.NewRequest("GET", "/groups/"+groupID.String(), nil)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result models.Group
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, expectedGroup.Name, result.Name)
	assert.Nil(t, result.Roles)
}
