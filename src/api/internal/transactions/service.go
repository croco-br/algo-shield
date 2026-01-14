package transactions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// Service defines the interface for transaction business logic
// This interface follows Dependency Inversion Principle
type Service interface {
	ProcessTransaction(ctx context.Context, event models.Event) (*asynq.TaskInfo, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
	ListTransactionsWithFilter(ctx context.Context, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error)
	ApproveTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	RejectTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
}

// TaskEnqueuer defines interface for enqueueing tasks to Asynq
// Follows Dependency Inversion Principle and Interface Segregation Principle
type TaskEnqueuer interface {
	EnqueueTransactionWithPriority(ctx context.Context, event models.Event, priority string) (*asynq.TaskInfo, error)
}

type service struct {
	repo     Repository
	enqueuer TaskEnqueuer
}

// NewService creates a new transaction service with dependency injection
// Follows Dependency Inversion Principle - receives interfaces, not concrete types
func NewService(repo Repository, enqueuer TaskEnqueuer) Service {
	return &service{
		repo:     repo,
		enqueuer: enqueuer,
	}
}

func (s *service) ProcessTransaction(ctx context.Context, event models.Event) (*asynq.TaskInfo, error) {
	// Determine priority based on event characteristics
	// This could be made more sophisticated with business rules
	priority := "default"

	// Example: Check if event has a priority field
	if p, ok := event["priority"].(string); ok {
		priority = p
	}

	// Enqueue transaction for processing
	return s.enqueuer.EnqueueTransactionWithPriority(ctx, event, priority)
}

func (s *service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *service) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(ctx, limit, offset)
}

func (s *service) ListTransactionsWithFilter(ctx context.Context, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error) {
	return s.repo.ListTransactionsWithFilter(ctx, filter, limit, offset)
}

func (s *service) ApproveTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.ApproveTransaction(ctx, id)
}

func (s *service) RejectTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.RejectTransaction(ctx, id)
}
