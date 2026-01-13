package transactions

import (
	"context"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
)

// RuleEvaluator defines the interface for rule evaluation
type RuleEvaluator interface {
	Evaluate(ctx context.Context, event models.Event) (*models.TransactionResult, error)
}

// Service handles transaction processing business logic
type Service struct {
	repo          Repository
	ruleEvaluator RuleEvaluator
}

// NewService creates a new transaction service with dependency injection
// Follows Dependency Inversion Principle - receives interface, not concrete type
func NewService(repo Repository, ruleEvaluator RuleEvaluator) *Service {
	return &Service{
		repo:          repo,
		ruleEvaluator: ruleEvaluator,
	}
}

// ProcessTransaction processes an event by evaluating rules and saving the result
func (s *Service) ProcessTransaction(ctx context.Context, event models.Event) error {
	// Evaluate event against rules
	result, err := s.ruleEvaluator.Evaluate(ctx, event)
	if err != nil {
		return err
	}

	// Create transaction record
	transactionID := uuid.New()
	now := time.Now()

	// Extract fields from generic event (with fallbacks for common field names)
	externalID := extractStringFromEvent(event, "external_id", "id", "event_id")
	amount := extractFloat64FromEvent(event, "amount", "value", "total")
	currency := extractStringFromEvent(event, "currency", "currency_code", "curr")
	origin := extractStringFromEvent(event, "origin", "from_account", "account", "user_id", "customer_id")
	destination := extractStringFromEvent(event, "destination", "to_account", "recipient_account", "recipient_id")
	eventType := extractStringFromEvent(event, "type", "transaction_type", "event_type")

	// Extract schema_id if present (added by synthetic event generator)
	var schemaID *uuid.UUID
	if schemaIDStr := extractStringFromEvent(event, "_schema_id"); schemaIDStr != "" {
		if parsed, err := uuid.Parse(schemaIDStr); err == nil {
			schemaID = &parsed
		}
		// Remove _schema_id from event to not include in metadata
		delete(event, "_schema_id")
	}

	// Check if this is a synthetic event
	isSynthetic := false
	if synth, ok := event["_synthetic"].(bool); ok && synth {
		isSynthetic = true
		// Remove _synthetic from event to not include in metadata
		delete(event, "_synthetic")
	}

	// Extract metadata if present
	var metadata map[string]any
	if meta, ok := event["metadata"].(map[string]any); ok {
		metadata = meta
	} else {
		// If no metadata field, use empty map
		metadata = make(map[string]any)
	}

	// Synthetic transactions always have pending status and are not processed
	var status models.TransactionStatus
	var processingTime int64
	var matchedRules []string
	var processedAt *time.Time
	if isSynthetic {
		status = models.StatusPending
		processingTime = 0
		matchedRules = []string{}
		processedAt = nil
	} else {
		status = result.Status
		processingTime = result.ProcessingTime
		matchedRules = result.MatchedRules
		processedAt = &now
	}

	transaction := &models.Transaction{
		ID:             transactionID,
		ExternalID:     externalID,
		SchemaID:       schemaID,
		Amount:         amount,
		Currency:       currency,
		Origin:         origin,
		Destination:    destination,
		Type:           eventType,
		Status:         status,
		ProcessingTime: processingTime,
		MatchedRules:   matchedRules,
		Metadata:       metadata,
		CreatedAt:      now,
		ProcessedAt:    processedAt,
	}

	// Save transaction to database (synthetic or regular table)
	var saveErr error
	if isSynthetic {
		saveErr = s.repo.SaveSyntheticTransaction(ctx, transaction)
	} else {
		saveErr = s.repo.SaveTransaction(ctx, transaction)
	}
	if saveErr != nil {
		return saveErr
	}

	log.Printf(
		"Processed transaction %s: status=%s, time=%dms",
		externalID, result.Status, result.ProcessingTime,
	)

	return nil
}

// Helper functions to extract values from generic event
func extractStringFromEvent(event models.Event, fieldNames ...string) string {
	for _, name := range fieldNames {
		if val, ok := event[name]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
	}
	return ""
}

func extractFloat64FromEvent(event models.Event, fieldNames ...string) float64 {
	for _, name := range fieldNames {
		if val, ok := event[name]; ok {
			if f, ok := toFloat64(val); ok {
				return f
			}
		}
	}
	return 0
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}
