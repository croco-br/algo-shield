package rules

import (
	"context"
	"time"

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

func (s *Service) CreateRule(ctx context.Context, rule *models.Rule) error {
	rule.ID = uuid.New()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	if err := s.repo.CreateRule(ctx, rule); err != nil {
		return err
	}

	// Invalidate cache for hot-reload
	s.redis.Del(ctx, "rules:cache")

	return nil
}

func (s *Service) GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error) {
	return s.repo.GetRule(ctx, id)
}

func (s *Service) ListRules(ctx context.Context) ([]models.Rule, error) {
	return s.repo.ListRules(ctx)
}

func (s *Service) UpdateRule(ctx context.Context, rule *models.Rule) error {
	rule.UpdatedAt = time.Now()

	if err := s.repo.UpdateRule(ctx, rule); err != nil {
		return err
	}

	// Invalidate cache for hot-reload
	s.redis.Del(ctx, "rules:cache")

	return nil
}

func (s *Service) DeleteRule(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteRule(ctx, id); err != nil {
		return err
	}

	// Invalidate cache for hot-reload
	s.redis.Del(ctx, "rules:cache")

	return nil
}
