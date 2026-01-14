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

	// Extract schema_id if present (added by synthetic event generator)
	var schemaID *uuid.UUID
	if schemaIDStr, ok := event["_schema_id"].(string); ok && schemaIDStr != "" {
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

	// The entire event becomes the metadata (all fields come from schema)
	metadata := make(map[string]any)
	for k, v := range event {
		metadata[k] = v
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
		SchemaID:       schemaID,
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
		transactionID, result.Status, result.ProcessingTime,
	)

	return nil
}
