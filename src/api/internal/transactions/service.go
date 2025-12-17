package transactions

import (
	"context"
	"encoding/json"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(db *pgxpool.Pool, redis *redis.Client) *Service {
	return &Service{
		repo:  NewPostgresRepository(db),
		redis: redis,
	}
}

func (s *Service) ProcessTransaction(ctx context.Context, event models.TransactionEvent) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Push to Redis queue
	return s.redis.LPush(ctx, "transaction:queue", eventJSON).Err()
}

func (s *Service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *Service) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	return s.repo.ListTransactions(ctx, limit, offset)
}
