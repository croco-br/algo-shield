package groups

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for group business logic
type Service interface {
	ListGroups(ctx context.Context) ([]models.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetGroupByName(ctx context.Context, name string) (*models.Group, error)
	LoadUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error)
	LoadGroupRoles(ctx context.Context, groupID uuid.UUID) ([]models.Role, error)
}

type service struct {
	repo Repository
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		repo: NewPostgresRepository(db),
	}
}

func (s *service) ListGroups(ctx context.Context) ([]models.Group, error) {
	return s.repo.ListGroups(ctx)
}

func (s *service) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	return s.repo.GetGroupByID(ctx, groupID)
}

func (s *service) GetGroupByName(ctx context.Context, name string) (*models.Group, error) {
	return s.repo.GetGroupByName(ctx, name)
}

func (s *service) LoadUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error) {
	return s.repo.LoadUserGroups(ctx, userID)
}

func (s *service) LoadGroupRoles(ctx context.Context, groupID uuid.UUID) ([]models.Role, error) {
	return s.repo.LoadGroupRoles(ctx, groupID)
}
