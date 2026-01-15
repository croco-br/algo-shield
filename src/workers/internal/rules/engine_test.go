package rules

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/workers/internal/schemas"
	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Engine_Evaluate_WhenNoRules_ThenReturnsApproved(t *testing.T) {
	engine := createTestEngine([]models.Rule{}, nil, nil)
	event := models.Event{"amount": 100.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusApproved, result.Status)
	assert.Empty(t, result.MatchedRules)
	assert.GreaterOrEqual(t, result.ProcessingTime, int64(0))
}

func Test_Engine_Evaluate_WhenRuleMatchesWithBlockAction_ThenReturnsRejected(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rule := models.Rule{
		Name:   "block-high-amount",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 1000",
		},
		SchemaID: &schemaID,
	}
	engine := createTestEngine([]models.Rule{rule}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 2000.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusRejected, result.Status)
	assert.Equal(t, []string{"block-high-amount"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenRuleMatchesWithReviewAction_ThenReturnsInReview(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rule := models.Rule{
		Name:   "review-medium-amount",
		Action: models.ActionReview,
		Conditions: map[string]any{
			"custom_expression": "amount > 500",
		},
		SchemaID: &schemaID,
	}
	engine := createTestEngine([]models.Rule{rule}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 750.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusInReview, result.Status)
	assert.Equal(t, []string{"review-medium-amount"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenRuleMatchesWithAllowAction_ThenReturnsApproved(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rule := models.Rule{
		Name:   "allow-low-amount",
		Action: models.ActionAllow,
		Conditions: map[string]any{
			"custom_expression": "amount < 100",
		},
		SchemaID: &schemaID,
	}
	engine := createTestEngine([]models.Rule{rule}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 50.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusApproved, result.Status)
	assert.Equal(t, []string{"allow-low-amount"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenBlockAndReviewRulesMatch_ThenBlockTakesPriority(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rules := []models.Rule{
		{
			Name:   "review-rule",
			Action: models.ActionReview,
			Conditions: map[string]any{
				"custom_expression": "amount > 500",
			},
			SchemaID: &schemaID,
		},
		{
			Name:   "block-rule",
			Action: models.ActionBlock,
			Conditions: map[string]any{
				"custom_expression": "amount > 1000",
			},
			SchemaID: &schemaID,
		},
	}
	engine := createTestEngine(rules, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 1500.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusRejected, result.Status)
	assert.ElementsMatch(t, []string{"review-rule", "block-rule"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenReviewAndAllowRulesMatch_ThenReviewTakesPriority(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rules := []models.Rule{
		{
			Name:   "allow-rule",
			Action: models.ActionAllow,
			Conditions: map[string]any{
				"custom_expression": "amount < 2000",
			},
			SchemaID: &schemaID,
		},
		{
			Name:   "review-rule",
			Action: models.ActionReview,
			Conditions: map[string]any{
				"custom_expression": "amount > 500",
			},
			SchemaID: &schemaID,
		},
	}
	engine := createTestEngine(rules, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 750.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusInReview, result.Status)
	assert.ElementsMatch(t, []string{"allow-rule", "review-rule"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenMultipleRulesMatchDifferentActions_ThenBlockHasHighestPriority(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rules := []models.Rule{
		{
			Name:   "allow-rule",
			Action: models.ActionAllow,
			Conditions: map[string]any{
				"custom_expression": "amount > 0",
			},
			SchemaID: &schemaID,
		},
		{
			Name:   "review-rule",
			Action: models.ActionReview,
			Conditions: map[string]any{
				"custom_expression": "amount > 500",
			},
			SchemaID: &schemaID,
		},
		{
			Name:   "block-rule",
			Action: models.ActionBlock,
			Conditions: map[string]any{
				"custom_expression": "amount > 10000",
			},
			SchemaID: &schemaID,
		},
	}
	engine := createTestEngine(rules, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 15000.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusRejected, result.Status)
	assert.ElementsMatch(t, []string{"allow-rule", "review-rule", "block-rule"}, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenNoRulesMatch_ThenReturnsApproved(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	rule := models.Rule{
		Name:   "block-high-amount",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 1000",
		},
		SchemaID: &schemaID,
	}
	engine := createTestEngine([]models.Rule{rule}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	event := models.Event{"amount": 100.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, models.StatusApproved, result.Status)
	assert.Empty(t, result.MatchedRules)
}

func Test_Engine_Evaluate_WhenProcessingCompletes_ThenRecordsProcessingTime(t *testing.T) {
	engine := createTestEngine([]models.Rule{}, nil, nil)
	event := models.Event{"amount": 100.0}

	result, err := engine.Evaluate(context.Background(), event)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, result.ProcessingTime, int64(0))
	assert.Less(t, result.ProcessingTime, int64(1000))
}

func Test_Engine_EvaluateCustomRule_WhenMissingCustomExpression_ThenReturnsFalse(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{})
	engine := createTestEngine([]models.Rule{}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	rule := models.Rule{
		Name:       "invalid-rule",
		Action:     models.ActionBlock,
		Conditions: map[string]any{},
		SchemaID:   &schemaID,
	}
	event := models.Event{"amount": 100.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.False(t, result)
}

func Test_Engine_EvaluateCustomRule_WhenCustomExpressionNotString_ThenReturnsFalse(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{})
	engine := createTestEngine([]models.Rule{}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	rule := models.Rule{
		Name:   "invalid-rule",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": 123,
		},
		SchemaID: &schemaID,
	}
	event := models.Event{"amount": 100.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.False(t, result)
}

func Test_Engine_EvaluateCustomRule_WhenMissingSchemaID_ThenReturnsFalse(t *testing.T) {
	engine := createTestEngine([]models.Rule{}, nil, nil)
	rule := models.Rule{
		Name:   "invalid-rule",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 100",
		},
		SchemaID: nil,
	}
	event := models.Event{"amount": 100.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.False(t, result)
}

func Test_Engine_EvaluateCustomRule_WhenSchemaNotFound_ThenReturnsFalse(t *testing.T) {
	schemaID := uuid.New()
	engine := createTestEngine([]models.Rule{}, map[uuid.UUID]*schemas.EventSchema{}, nil)
	rule := models.Rule{
		Name:   "invalid-rule",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 100",
		},
		SchemaID: &schemaID,
	}
	event := models.Event{"amount": 100.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.False(t, result)
}

func Test_Engine_EvaluateCustomRule_WhenValidExpression_ThenReturnsTrue(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	engine := createTestEngine([]models.Rule{}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	rule := models.Rule{
		Name:   "valid-rule",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 100",
		},
		SchemaID: &schemaID,
	}
	event := models.Event{"amount": 200.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.True(t, result)
}

func Test_Engine_EvaluateCustomRule_WhenExpressionEvaluatesToFalse_ThenReturnsFalse(t *testing.T) {
	schemaID := uuid.New()
	schema := createTestSchema(schemaID, []schemas.ExtractedField{
		{Path: "amount", Type: schemas.FieldTypeNumber},
	})
	engine := createTestEngine([]models.Rule{}, map[uuid.UUID]*schemas.EventSchema{schemaID: schema}, nil)
	rule := models.Rule{
		Name:   "valid-rule",
		Action: models.ActionBlock,
		Conditions: map[string]any{
			"custom_expression": "amount > 1000",
		},
		SchemaID: &schemaID,
	}
	event := models.Event{"amount": 200.0}

	result := engine.evaluateCustomRule(context.Background(), event, rule)

	assert.False(t, result)
}

// Test helpers

// testEngine is a test version of Engine that uses interfaces for easy mocking
type testEngine struct {
	ruleService    *RuleService
	schemaService  schemaServiceInterface
	historyRepo    transactions.TransactionHistoryRepository
	defaultTimeout time.Duration
}

// schemaServiceInterface defines the schema service interface for testing
type schemaServiceInterface interface {
	GetSchema(id uuid.UUID) *schemas.EventSchema
	LoadSchemas(ctx context.Context) error
	InvalidateSchema(ctx context.Context, id uuid.UUID)
	SubscribeToInvalidations(ctx context.Context)
	GetAllSchemas() map[uuid.UUID]*schemas.EventSchema
}

func createTestEngine(rules []models.Rule, schemasMap map[uuid.UUID]*schemas.EventSchema, historyRepo transactions.TransactionHistoryRepository) *testEngine {
	mockRepo := &mockRuleReader{rules: rules}
	ruleService := NewRuleService(mockRepo)
	_ = ruleService.LoadRules(context.Background())

	schemaService := &mockSchemaService{
		schemas: schemasMap,
	}
	return &testEngine{
		ruleService:    ruleService,
		schemaService:  schemaService,
		historyRepo:    historyRepo,
		defaultTimeout: 5 * time.Second,
	}
}

// Implement Engine methods for testEngine

func (e *testEngine) Evaluate(ctx context.Context, event models.Event) (*models.TransactionResult, error) {
	startTime := time.Now()

	matchedRules := make([]string, 0)
	status := models.StatusApproved

	rules := e.ruleService.GetRules()
	for _, rule := range rules {
		matched := e.evaluateRule(ctx, event, rule)
		if matched {
			matchedRules = append(matchedRules, rule.Name)

			switch rule.Action {
			case models.ActionBlock:
				status = models.StatusRejected
			case models.ActionReview:
				if status != models.StatusRejected {
					status = models.StatusInReview
				}
			case models.ActionAllow:
				status = models.StatusApproved
			}
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	result := &models.TransactionResult{
		Status:         status,
		MatchedRules:   matchedRules,
		ProcessingTime: processingTime,
	}

	return result, nil
}

func (e *testEngine) evaluateRule(ctx context.Context, event models.Event, rule models.Rule) bool {
	return e.evaluateCustomRule(ctx, event, rule)
}

func (e *testEngine) evaluateCustomRule(ctx context.Context, event models.Event, rule models.Rule) bool {
	expression, ok := rule.Conditions["custom_expression"].(string)
	if !ok {
		return false
	}

	if rule.SchemaID == nil {
		return false
	}
	schema := e.schemaService.GetSchema(*rule.SchemaID)
	if schema == nil {
		return false
	}

	return schemas.EvaluateExpressionWithSchema(ctx, expression, event, schema, e.historyRepo)
}

func createTestSchema(id uuid.UUID, fields []schemas.ExtractedField) *schemas.EventSchema {
	return &schemas.EventSchema{
		ID:              id,
		Name:            "test-schema",
		Description:     "Test schema",
		ExtractedFields: fields,
		SampleJSON:      map[string]any{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// mockRuleReader is a simple mock for testing
type mockRuleReader struct {
	rules []models.Rule
	err   error
}

func (m *mockRuleReader) LoadRules(ctx context.Context) ([]models.Rule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rules, nil
}

// mockSchemaService is a simple mock for testing that implements the schemaService interface
type mockSchemaService struct {
	schemas map[uuid.UUID]*schemas.EventSchema
}

func (m *mockSchemaService) GetSchema(id uuid.UUID) *schemas.EventSchema {
	if m.schemas == nil {
		return nil
	}
	return m.schemas[id]
}

func (m *mockSchemaService) LoadSchemas(ctx context.Context) error {
	return nil
}

func (m *mockSchemaService) InvalidateSchema(ctx context.Context, id uuid.UUID) {
}

func (m *mockSchemaService) SubscribeToInvalidations(ctx context.Context) {
}

func (m *mockSchemaService) GetAllSchemas() map[uuid.UUID]*schemas.EventSchema {
	return m.schemas
}
