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
	Evaluate(ctx context.Context, event models.TransactionEvent) (*models.TransactionResult, error)
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

// ProcessTransaction processes a transaction event by evaluating rules and saving the result
func (s *Service) ProcessTransaction(ctx context.Context, event models.TransactionEvent) error {
	// Evaluate transaction against rules
	result, err := s.ruleEvaluator.Evaluate(ctx, event)
	if err != nil {
		return err
	}

	// Create transaction record
	transactionID := uuid.New()
	now := time.Now()

	transaction := &models.Transaction{
		ID:             transactionID,
		ExternalID:     event.ExternalID,
		Amount:         event.Amount,
		Currency:       event.Currency,
		FromAccount:    event.FromAccount,
		ToAccount:      event.ToAccount,
		Type:           event.Type,
		Status:         result.Status,
		RiskScore:      result.RiskScore,
		RiskLevel:      result.RiskLevel,
		ProcessingTime: result.ProcessingTime, // Use ProcessingTime from result
		MatchedRules:   result.MatchedRules,
		Metadata:       event.Metadata,
		CreatedAt:      now,
		ProcessedAt:    &now,
	}

	// Save transaction to database
	if err := s.repo.SaveTransaction(ctx, transaction); err != nil {
		return err
	}

	log.Printf(
		"Processed transaction %s: status=%s, risk_score=%.2f, risk_level=%s, time=%dms",
		event.ExternalID, result.Status, result.RiskScore, result.RiskLevel, result.ProcessingTime,
	)

	return nil
}
