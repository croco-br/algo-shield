package rules

import (
	"context"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// TransactionHistoryProvider defines the interface for transaction history queries
type TransactionHistoryProvider interface {
	CountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (int, error)
}

// Engine evaluates transactions against rules
type Engine struct {
	ruleService     *RuleService
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

	return &Engine{
		ruleService:     ruleService,
		historyProvider: historyProvider,
		defaultTimeout:  ruleEvaluationTimeout,
	}
}

// LoadRules loads rules from the rule loader
func (e *Engine) LoadRules(ctx context.Context) error {
	return e.ruleService.LoadRules(ctx)
}

// Evaluate evaluates a transaction against all loaded rules
func (e *Engine) Evaluate(ctx context.Context, event models.TransactionEvent) (*models.TransactionResult, error) {
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

// evaluateRule evaluates a single rule against a transaction
func (e *Engine) evaluateRule(ctx context.Context, event models.TransactionEvent, rule models.Rule) bool {
	switch rule.Type {
	case models.RuleTypeAmount:
		return e.evaluateAmountRule(event, rule)
	case models.RuleTypeVelocity:
		return e.evaluateVelocityRule(ctx, event, rule)
	case models.RuleTypeBlocklist:
		return e.evaluateBlocklistRule(event, rule)
	case models.RuleTypeGeography:
		return e.evaluateGeographyRule(event, rule)
	case models.RuleTypeCustom:
		return e.evaluateCustomRule(event, rule)
	default:
		log.Printf("Unknown rule type: %s", rule.Type)
		return false
	}
}

// evaluateAmountRule checks if transaction amount exceeds threshold
func (e *Engine) evaluateAmountRule(event models.TransactionEvent, rule models.Rule) bool {
	threshold, ok := rule.Conditions["amount_threshold"].(float64)
	if !ok {
		log.Printf("Amount rule missing or invalid amount_threshold condition")
		return false
	}
	return event.Amount > threshold
}

// evaluateVelocityRule checks transaction velocity (count in time window)
func (e *Engine) evaluateVelocityRule(ctx context.Context, event models.TransactionEvent, rule models.Rule) bool {
	timeWindow, ok := rule.Conditions["time_window_seconds"].(float64)
	if !ok {
		timeWindow = 3600 // default 1 hour
	}

	maxCount, ok := rule.Conditions["transaction_count"].(float64)
	if !ok {
		log.Printf("Velocity rule missing transaction_count condition")
		return false
	}

	// Use history provider instead of direct DB access
	ctx, cancel := context.WithTimeout(ctx, e.defaultTimeout)
	defer cancel()

	count, err := e.historyProvider.CountByAccountInTimeWindow(ctx, event.FromAccount, int(timeWindow))
	if err != nil {
		log.Printf("Failed to query transaction history for velocity rule: %v", err)
		return false
	}

	return count > int(maxCount)
}

// evaluateBlocklistRule checks if accounts are in blocklist
func (e *Engine) evaluateBlocklistRule(event models.TransactionEvent, rule models.Rule) bool {
	blocklist, ok := rule.Conditions["blocklisted_accounts"].([]interface{})
	if !ok {
		log.Printf("Blocklist rule missing or invalid blocklisted_accounts condition")
		return false
	}

	for _, account := range blocklist {
		if accountStr, ok := account.(string); ok {
			if accountStr == event.FromAccount || accountStr == event.ToAccount {
				return true
			}
		}
	}

	return false
}

// evaluateGeographyRule checks if transaction coordinates are within restricted polygon
func (e *Engine) evaluateGeographyRule(event models.TransactionEvent, rule models.Rule) bool {
	polygonStr, ok := rule.Conditions["polygon_coordinates"].(string)
	if !ok {
		log.Printf("Geography rule missing or invalid polygon_coordinates condition")
		return false
	}

	polygon, err := ParsePolygon(polygonStr)
	if err != nil {
		log.Printf("Geography rule: failed to parse polygon coordinates: %v", err)
		return false
	}

	if len(polygon) < 3 {
		log.Printf("Geography rule: polygon must have at least 3 points")
		return false
	}

	lat, lon, ok := GetCoordinatesFromMetadata(event.Metadata)
	if !ok {
		// No coordinates in transaction metadata, rule doesn't match
		return false
	}

	return PointInPolygon(lat, lon, polygon)
}

// evaluateCustomRule evaluates custom expression against transaction
func (e *Engine) evaluateCustomRule(event models.TransactionEvent, rule models.Rule) bool {
	expression, ok := rule.Conditions["custom_expression"].(string)
	if !ok {
		log.Printf("Custom rule missing or invalid custom_expression condition")
		return false
	}

	return EvaluateExpression(expression, event)
}
