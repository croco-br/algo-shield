package transactions

import (
	"context"
	"encoding/json"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Service defines the interface for transaction business logic
// This interface follows Dependency Inversion Principle
type Service interface {
	ProcessTransaction(ctx context.Context, event models.TransactionEvent) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
}

// QueuePusher defines interface for pushing to queue
// Follows Dependency Inversion Principle
type QueuePusher interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type service struct {
	repo      Repository
	queuePush QueuePusher
}

// NewService creates a new transaction service with dependency injection
// Follows Dependency Inversion Principle - receives interfaces, not concrete types
func NewService(repo Repository, queuePush QueuePusher) Service {
	return &service{
		repo:      repo,
		queuePush: queuePush,
	}
}

func (s *service) ProcessTransaction(ctx context.Context, event models.TransactionEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Push to Redis queue
	return s.queuePush.LPush(ctx, "transaction:queue", eventJSON).Err()
}

func (s *service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *service) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(ctx, limit, offset)
}
