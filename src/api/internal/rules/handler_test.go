package rules

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Handler_NewHandler_WhenCalled_ThenReturnsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)

	handler := NewHandler(repo)

	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func Test_Handler_CreateRule_WhenValidRule_ThenCreatesRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Post("/rules", handler.CreateRule)

	rule := models.Rule{
		Name:        "Test Rule",
		Description: "Test Description",
		Action:      models.ActionBlock,
		Priority:    50,
		Enabled:     true,
		Conditions:  map[string]any{"amount": ">1000"},
	}

	repo.EXPECT().CreateRule(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, r *models.Rule) error {
			assert.Equal(t, rule.Name, r.Name)
			assert.NotEqual(t, uuid.Nil, r.ID)
			assert.False(t, r.CreatedAt.IsZero())
			assert.False(t, r.UpdatedAt.IsZero())
			return nil
		},
	)

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func Test_Handler_CreateRule_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Post("/rules", handler.CreateRule)

	req := httptest.NewRequest("POST", "/rules", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Invalid request body")
}

func Test_Handler_CreateRule_WhenValidationFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Post("/rules", handler.CreateRule)

	rule := models.Rule{
		Name:        "",
		Description: "Test",
		Action:      models.ActionBlock,
		Priority:    50,
		Enabled:     true,
		Conditions:  map[string]any{"amount": ">1000"},
	}

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_CreateRule_WhenRepositoryFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Post("/rules", handler.CreateRule)

	rule := models.Rule{
		Name:        "Test Rule",
		Description: "Test Description",
		Action:      models.ActionBlock,
		Priority:    50,
		Enabled:     true,
		Conditions:  map[string]any{"amount": ">1000"},
	}

	repo.EXPECT().CreateRule(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("POST", "/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(respBody), "Failed to create rule")
}

func Test_Handler_GetRule_WhenRuleExists_ThenReturnsRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules/:id", handler.GetRule)

	ruleID := uuid.New()
	expectedRule := &models.Rule{
		ID:          ruleID,
		Name:        "Test Rule",
		Description: "Test Description",
		Action:      models.ActionBlock,
		Priority:    50,
		Enabled:     true,
		Conditions:  map[string]any{"amount": ">1000"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.EXPECT().GetRule(gomock.Any(), ruleID).Return(expectedRule, nil)

	req := httptest.NewRequest("GET", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result models.Rule
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, expectedRule.ID, result.ID)
	assert.Equal(t, expectedRule.Name, result.Name)
}

func Test_Handler_GetRule_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules/:id", handler.GetRule)

	req := httptest.NewRequest("GET", "/rules/invalid-uuid", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Invalid rule ID")
}

func Test_Handler_GetRule_WhenRuleNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules/:id", handler.GetRule)

	ruleID := uuid.New()
	repo.EXPECT().GetRule(gomock.Any(), ruleID).Return(nil, pgx.ErrNoRows)

	req := httptest.NewRequest("GET", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Rule not found")
}

func Test_Handler_GetRule_WhenRepositoryFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules/:id", handler.GetRule)

	ruleID := uuid.New()
	repo.EXPECT().GetRule(gomock.Any(), ruleID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Failed to fetch rule")
}

func Test_Handler_ListRules_WhenRulesExist_ThenReturnsRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules", handler.ListRules)

	expectedRules := []models.Rule{
		{
			ID:          uuid.New(),
			Name:        "Rule 1",
			Description: "Description 1",
			Action:      models.ActionBlock,
			Priority:    10,
			Enabled:     true,
			Conditions:  map[string]any{"amount": ">1000"},
		},
		{
			ID:          uuid.New(),
			Name:        "Rule 2",
			Description: "Description 2",
			Action:      models.ActionReview,
			Priority:    20,
			Enabled:     true,
			Conditions:  map[string]any{"amount": ">5000"},
		},
	}

	repo.EXPECT().ListRules(gomock.Any()).Return(expectedRules, nil)

	req := httptest.NewRequest("GET", "/rules", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string][]models.Rule
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Len(t, result["rules"], 2)
	assert.Equal(t, expectedRules[0].Name, result["rules"][0].Name)
}

func Test_Handler_ListRules_WhenRepositoryFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Get("/rules", handler.ListRules)

	repo.EXPECT().ListRules(gomock.Any()).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/rules", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Failed to fetch rules")
}

func Test_Handler_UpdateRule_WhenValidRule_ThenUpdatesRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	ruleID := uuid.New()
	rule := models.Rule{
		Name:        "Updated Rule",
		Description: "Updated Description",
		Action:      models.ActionAllow,
		Priority:    75,
		Enabled:     false,
		Conditions:  map[string]any{"amount": "<500"},
	}

	repo.EXPECT().UpdateRule(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, r *models.Rule) error {
			assert.Equal(t, ruleID, r.ID)
			assert.Equal(t, rule.Name, r.Name)
			assert.False(t, r.UpdatedAt.IsZero())
			return nil
		},
	)

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("PUT", "/rules/"+ruleID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func Test_Handler_UpdateRule_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	rule := models.Rule{
		Name:        "Updated Rule",
		Description: "Updated Description",
		Action:      models.ActionAllow,
		Priority:    75,
		Enabled:     false,
		Conditions:  map[string]any{"amount": "<500"},
	}

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("PUT", "/rules/invalid-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(respBody), "Invalid rule ID")
}

func Test_Handler_UpdateRule_WhenInvalidJSON_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	ruleID := uuid.New()

	req := httptest.NewRequest("PUT", "/rules/"+ruleID.String(), bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Invalid request body")
}

func Test_Handler_UpdateRule_WhenValidationFails_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	ruleID := uuid.New()
	rule := models.Rule{
		Name:        "",
		Description: "Updated Description",
		Action:      models.ActionAllow,
		Priority:    75,
		Enabled:     false,
		Conditions:  map[string]any{"amount": "<500"},
	}

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("PUT", "/rules/"+ruleID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func Test_Handler_UpdateRule_WhenRuleNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	ruleID := uuid.New()
	rule := models.Rule{
		Name:        "Updated Rule",
		Description: "Updated Description",
		Action:      models.ActionAllow,
		Priority:    75,
		Enabled:     false,
		Conditions:  map[string]any{"amount": "<500"},
	}

	repo.EXPECT().UpdateRule(gomock.Any(), gomock.Any()).Return(pgx.ErrNoRows)

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("PUT", "/rules/"+ruleID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(respBody), "Rule not found")
}

func Test_Handler_UpdateRule_WhenRepositoryFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Put("/rules/:id", handler.UpdateRule)

	ruleID := uuid.New()
	rule := models.Rule{
		Name:        "Updated Rule",
		Description: "Updated Description",
		Action:      models.ActionAllow,
		Priority:    75,
		Enabled:     false,
		Conditions:  map[string]any{"amount": "<500"},
	}

	repo.EXPECT().UpdateRule(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

	body, err := json.Marshal(rule)
	require.NoError(t, err)
	req := httptest.NewRequest("PUT", "/rules/"+ruleID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(respBody), "Failed to update rule")
}

func Test_Handler_DeleteRule_WhenRuleExists_ThenDeletesRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Delete("/rules/:id", handler.DeleteRule)

	ruleID := uuid.New()
	repo.EXPECT().DeleteRule(gomock.Any(), ruleID).Return(nil)

	req := httptest.NewRequest("DELETE", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func Test_Handler_DeleteRule_WhenInvalidID_ThenReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Delete("/rules/:id", handler.DeleteRule)

	req := httptest.NewRequest("DELETE", "/rules/invalid-uuid", nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Invalid rule ID")
}

func Test_Handler_DeleteRule_WhenRuleNotFound_ThenReturnsNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Delete("/rules/:id", handler.DeleteRule)

	ruleID := uuid.New()
	repo.EXPECT().DeleteRule(gomock.Any(), ruleID).Return(pgx.ErrNoRows)

	req := httptest.NewRequest("DELETE", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Rule not found")
}

func Test_Handler_DeleteRule_WhenRepositoryFails_ThenReturnsInternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockRepository(ctrl)
	handler := NewHandler(repo)

	app := fiber.New()
	app.Delete("/rules/:id", handler.DeleteRule)

	ruleID := uuid.New()
	repo.EXPECT().DeleteRule(gomock.Any(), ruleID).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/rules/"+ruleID.String(), nil)

	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "Failed to delete rule")
}
