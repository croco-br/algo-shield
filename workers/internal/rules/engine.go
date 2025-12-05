package rules

import (
	"context"
	"encoding/json"
	"time"

	"github.com/algo-shield/algo-shield/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Engine struct {
	db    *pgxpool.Pool
	redis *redis.Client
	rules []models.Rule
}

func NewEngine(db *pgxpool.Pool, redis *redis.Client) *Engine {
	return &Engine{
		db:    db,
		redis: redis,
		rules: make([]models.Rule, 0),
	}
}

// LoadRules loads rules from database or cache
func (e *Engine) LoadRules(ctx context.Context) error {
	// Try to get from cache first
	cachedRules, err := e.redis.Get(ctx, "rules:cache").Result()
	if err == nil && cachedRules != "" {
		return json.Unmarshal([]byte(cachedRules), &e.rules)
	}

	// Load from database
	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		WHERE enabled = true
		ORDER BY priority ASC
	`

	rows, err := e.db.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	rules := make([]models.Rule, 0)
	for rows.Next() {
		var rule models.Rule
		var conditionsJSON []byte

		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Action,
			&rule.Priority, &rule.Enabled, &conditionsJSON, &rule.Score,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(conditionsJSON, &rule.Conditions)
		rules = append(rules, rule)
	}

	e.rules = rules

	// Cache rules
	rulesJSON, _ := json.Marshal(rules)
	e.redis.Set(ctx, "rules:cache", rulesJSON, 5*time.Minute)

	return nil
}

// Evaluate evaluates a transaction against all loaded rules
func (e *Engine) Evaluate(ctx context.Context, event models.TransactionEvent) (*models.TransactionResult, error) {
	startTime := time.Now()

	riskScore := 0.0
	matchedRules := make([]string, 0)
	status := models.StatusApproved
	riskLevel := models.RiskLow

	// Evaluate each rule
	for _, rule := range e.rules {
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
	case models.RuleTypeBlacklist:
		return e.evaluateBlacklistRule(event, rule)
	case models.RuleTypePattern:
		return e.evaluatePatternRule(event, rule)
	default:
		return false
	}
}

func (e *Engine) evaluateAmountRule(event models.TransactionEvent, rule models.Rule) bool {
	threshold, ok := rule.Conditions["amount_threshold"].(float64)
	if !ok {
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
		return false
	}

	// Query recent transactions from same account
	query := `
		SELECT COUNT(*) 
		FROM transactions 
		WHERE from_account = $1 
		AND created_at > NOW() - INTERVAL '1 second' * $2
	`

	var count int
	err := e.db.QueryRow(ctx, query, event.FromAccount, int(timeWindow)).Scan(&count)
	if err != nil {
		return false
	}

	return count > int(maxCount)
}

func (e *Engine) evaluateBlacklistRule(event models.TransactionEvent, rule models.Rule) bool {
	blacklist, ok := rule.Conditions["blacklisted_accounts"].([]interface{})
	if !ok {
		return false
	}

	for _, account := range blacklist {
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
		return false
	}

	// Example: Check if transaction type matches pattern
	return event.Type == pattern
}

