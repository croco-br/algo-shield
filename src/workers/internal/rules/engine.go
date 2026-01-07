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

// Engine evaluates events against rules using schemas
type Engine struct {
	ruleService    *RuleService
	schemaService  *schemas.SchemaService
	historyRepo    transactions.TransactionHistoryRepository
	defaultTimeout time.Duration
}

// NewEngine creates a new rule engine
func NewEngine(db *pgxpool.Pool, redis *redis.Client, ruleEvaluationTimeout time.Duration) *Engine {
	// Create rule repository and service with dependency injection
	ruleRepo := rules.NewPostgresRepository(db, redis)
	ruleService := NewRuleService(ruleRepo)

	// Create schema repository and service
	schemaRepo := schemas.NewPostgresRepository(db)
	schemaService := schemas.NewSchemaService(schemaRepo, redis)

	// Create history repository for velocity helpers
	historyRepo := transactions.NewPostgresHistoryRepository(db)

	return &Engine{
		ruleService:    ruleService,
		schemaService:  schemaService,
		historyRepo:    historyRepo,
		defaultTimeout: ruleEvaluationTimeout,
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
// This is a blocking function that should be called in a goroutine managed by errgroup
func (e *Engine) StartSchemaInvalidationSubscription(ctx context.Context) {
	e.schemaService.SubscribeToInvalidations(ctx)
}

// Evaluate evaluates an event against all loaded rules using schema-based evaluation
func (e *Engine) Evaluate(ctx context.Context, event models.Event) (*models.TransactionResult, error) {
	startTime := time.Now()

	matchedRules := make([]string, 0)
	status := models.StatusApproved

	// Evaluate each rule
	rules := e.ruleService.GetRules()
	for _, rule := range rules {
		matched := e.evaluateRule(ctx, event, rule)
		if matched {
			matchedRules = append(matchedRules, rule.Name)

			// Determine action
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

// evaluateRule evaluates a single rule against an event
// All rules use custom expressions (schema-based)
func (e *Engine) evaluateRule(ctx context.Context, event models.Event, rule models.Rule) bool {
	return e.evaluateCustomRule(ctx, event, rule)
}

// evaluateCustomRule evaluates custom expression against event using schema
func (e *Engine) evaluateCustomRule(ctx context.Context, event models.Event, rule models.Rule) bool {
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

	// Use schema-based expression evaluation with history repository for velocity helpers
	return schemas.EvaluateExpressionWithSchema(ctx, expression, event, schema, e.historyRepo)
}
