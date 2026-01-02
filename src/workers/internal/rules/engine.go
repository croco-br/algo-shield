package rules

import (
	"context"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/algo-shield/algo-shield/src/workers/internal/schemas"
	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TransactionHistoryProvider defines the interface for transaction history queries
type TransactionHistoryProvider interface {
	CountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (int, error)
}

// Engine evaluates events against rules using schemas
type Engine struct {
	ruleService     *RuleService
	schemaService   *schemas.SchemaService
	historyProvider TransactionHistoryProvider
	defaultTimeout  time.Duration
}

// NewEngine creates a new rule engine
func NewEngine(db *pgxpool.Pool, redis *redis.Client, historyProvider TransactionHistoryProvider, ruleEvaluationTimeout time.Duration) *Engine {
	if historyProvider == nil {
		// Default to PostgresHistoryRepository if not provided
		historyProvider = transactions.NewPostgresHistoryRepository(db)
	}

	// Create rule repository and service with dependency injection
	ruleRepo := rules.NewPostgresRepository(db, redis)
	ruleService := NewRuleService(ruleRepo)

	// Create schema repository and service
	schemaRepo := schemas.NewPostgresRepository(db)
	schemaService := schemas.NewSchemaService(schemaRepo, redis)

	return &Engine{
		ruleService:     ruleService,
		schemaService:   schemaService,
		historyProvider: historyProvider,
		defaultTimeout:  ruleEvaluationTimeout,
	}
}

// LoadRules loads rules and schemas
func (e *Engine) LoadRules(ctx context.Context) error {
	if err := e.ruleService.LoadRules(ctx); err != nil {
		return err
	}
	return e.schemaService.LoadSchemas(ctx)
}

// StartSchemaInvalidationSubscription starts listening for schema changes
func (e *Engine) StartSchemaInvalidationSubscription(ctx context.Context) {
	go e.schemaService.SubscribeToInvalidations(ctx)
}

// Evaluate evaluates an event against all loaded rules using schema-based evaluation
func (e *Engine) Evaluate(ctx context.Context, event models.Event) (*models.TransactionResult, error) {
	startTime := time.Now()

	riskScore := 0.0
	matchedRules := make([]string, 0)
	status := models.StatusApproved

	// Evaluate each rule
	rules := e.ruleService.GetRules()
	for _, rule := range rules {
		matched := e.evaluateRule(ctx, event, rule)
		if matched {
			matchedRules = append(matchedRules, rule.Name)
			riskScore += rule.Score

			// Determine action
			switch rule.Action {
			case models.ActionBlock:
				status = models.StatusRejected
			case models.ActionReview:
				if status != models.StatusRejected {
					status = models.StatusReview
				}
			case models.ActionScore:
				// Just add to score
			}
		}
	}

	// Determine risk level based on score
	var riskLevel models.RiskLevel
	switch {
	case riskScore >= 80:
		riskLevel = models.RiskHigh
		if status == models.StatusApproved {
			status = models.StatusReview
		}
	case riskScore >= 50:
		riskLevel = models.RiskMedium
	default:
		riskLevel = models.RiskLow
	}

	processingTime := time.Since(startTime).Milliseconds()

	result := &models.TransactionResult{
		Status:         status,
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		MatchedRules:   matchedRules,
		ProcessingTime: processingTime,
	}

	return result, nil
}

// evaluateRule evaluates a single rule against an event
// Only custom rules are supported (schema-based expressions)
func (e *Engine) evaluateRule(ctx context.Context, event models.Event, rule models.Rule) bool {
	if rule.Type != models.RuleTypeCustom {
		log.Printf("Unsupported rule type: %s (only 'custom' is supported)", rule.Type)
		return false
	}
	return e.evaluateCustomRule(event, rule)
}

// evaluateCustomRule evaluates custom expression against event using schema
func (e *Engine) evaluateCustomRule(event models.Event, rule models.Rule) bool {
	expression, ok := rule.Conditions["custom_expression"].(string)
	if !ok {
		log.Printf("Custom rule missing or invalid custom_expression condition")
		return false
	}

	// Get schema for this rule
	if rule.SchemaID == nil {
		log.Printf("Custom rule missing schema_id")
		return false
	}
	schema := e.schemaService.GetSchema(*rule.SchemaID)
	if schema == nil {
		log.Printf("Schema %s not found for rule %s", *rule.SchemaID, rule.Name)
		return false
	}

	// Use schema-based expression evaluation
	return schemas.EvaluateExpressionWithSchema(expression, event, schema)
}
