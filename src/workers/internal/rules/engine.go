package engine

import (
	"context"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
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
	ruleLoader      *RuleLoader
	historyProvider TransactionHistoryProvider
	defaultTimeout  time.Duration
}

// NewEngine creates a new rule engine
func NewEngine(db *pgxpool.Pool, redis *redis.Client, historyProvider TransactionHistoryProvider) *Engine {
	if historyProvider == nil {
		// Default to PostgresHistoryRepository if not provided
		historyProvider = transactions.NewPostgresHistoryRepository(db)
	}

	return &Engine{
		ruleLoader:      NewRuleLoader(db, redis),
		historyProvider: historyProvider,
		defaultTimeout:  300 * time.Millisecond,
	}
}

// LoadRules loads rules from the rule loader
func (e *Engine) LoadRules(ctx context.Context) error {
	return e.ruleLoader.LoadRules(ctx)
}

// Evaluate evaluates a transaction against all loaded rules
func (e *Engine) Evaluate(ctx context.Context, event models.TransactionEvent) (*models.TransactionResult, error) {
	startTime := time.Now()

	riskScore := 0.0
	matchedRules := make([]string, 0)
	status := models.StatusApproved

	// Evaluate each rule
	rules := e.ruleLoader.GetRules()
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
	case models.RuleTypePattern:
		return e.evaluatePatternRule(event, rule)
	case models.RuleTypeGeography:
		return e.evaluateGeographyRule(event, rule)
	case models.RuleTypeCustom:
		return e.evaluateCustomRule(ctx, event, rule)
	default:
		log.Printf("Unknown rule type: %s", rule.Type)
		return false
	}
}

func (e *Engine) evaluateAmountRule(event models.TransactionEvent, rule models.Rule) bool {
	threshold, ok := rule.Conditions["amount_threshold"].(float64)
	if !ok {
		log.Printf("Amount rule missing or invalid amount_threshold condition")
		return false
	}
	return event.Amount > threshold
}

func (e *Engine) evaluateVelocityRule(ctx context.Context, event models.TransactionEvent, rule models.Rule) bool {
	// Check transaction velocity (count in time window)
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

func (e *Engine) evaluatePatternRule(event models.TransactionEvent, rule models.Rule) bool {
	// Simple pattern matching - can be extended
	pattern, ok := rule.Conditions["pattern"].(string)
	if !ok {
		log.Printf("Pattern rule missing or invalid pattern condition")
		return false
	}

	// Example: Check if transaction type matches pattern
	return event.Type == pattern
}

func (e *Engine) evaluateGeographyRule(event models.TransactionEvent, rule models.Rule) bool {
	// Check if transaction involves restricted geographic regions
	restrictedRegions, ok := rule.Conditions["restricted_regions"].([]interface{})
	if !ok {
		log.Printf("Geography rule missing or invalid restricted_regions condition")
		return false
	}

	// Get transaction region from metadata
	region, ok := event.Metadata["region"].(string)
	if !ok {
		// If no region in metadata, check country
		country, ok := event.Metadata["country"].(string)
		if !ok {
			return false
		}
		region = country
	}

	// Check if region is in restricted list
	for _, restricted := range restrictedRegions {
		if restrictedStr, ok := restricted.(string); ok {
			if restrictedStr == region {
				return true
			}
		}
	}

	return false
}

func (e *Engine) evaluateCustomRule(ctx context.Context, event models.TransactionEvent, rule models.Rule) bool {
	// Custom expression evaluation
	expression, ok := rule.Conditions["custom_expression"].(string)
	if !ok {
		log.Printf("Custom rule missing or invalid custom_expression condition")
		return false
	}

	if expression == "" {
		return false
	}

	// Basic custom rule evaluation
	// Supports simple field checks from metadata or event fields
	// Format examples:
	// - "amount > 1000"
	// - "currency == USD"
	// - "type == transfer"
	// - "metadata.key == value"

	// For now, evaluate based on metadata conditions
	// If custom_expression is a key in metadata, check if it evaluates to true
	if value, exists := event.Metadata[expression]; exists {
		// If it's a boolean, return it
		if boolVal, ok := value.(bool); ok {
			return boolVal
		}
		// If it's a string "true", return true
		if strVal, ok := value.(string); ok && strVal == "true" {
			return true
		}
	}

	// Check if expression matches event fields
	// Simple pattern: "field:value" format
	// This is a basic implementation - for complex expressions, use a proper expression engine
	// like github.com/antonmedv/expr

	log.Printf("Custom rule expression evaluation: %s (basic implementation)", expression)
	// Return false for now - requires proper expression parser for full implementation
	return false
}
