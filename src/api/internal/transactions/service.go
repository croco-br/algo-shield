package transactions

import (
	"context"
	"encoding/json"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Service defines the interface for transaction business logic
// This interface follows Dependency Inversion Principle
type Service interface {
	ProcessTransaction(ctx context.Context, event models.TransactionEvent) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
}

type service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(db *pgxpool.Pool, redis *redis.Client) Service {
	return &service{
		repo:  NewPostgresRepository(db),
		redis: redis,
	}
}

func (s *service) ProcessTransaction(ctx context.Context, event models.TransactionEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Push to Redis queue
	return s.redis.LPush(ctx, "transaction:queue", eventJSON).Err()
}

func (s *service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *service) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(ctx, limit, offset)
}
