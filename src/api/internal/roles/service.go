package roles

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service defines the interface for role business logic
type Service interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error
	LoadUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)
}

// RepositoryAccess provides access to repository methods that need transactions
// This allows other services (like user) to use role repository methods with transactions
type RepositoryAccess interface {
	GetRoleIDByName(ctx context.Context, tx interface{}, name string) (uuid.UUID, error)
	AssignRoleToUser(ctx context.Context, tx interface{}, userID, roleID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		repo: NewPostgresRepository(db),
	}
}

func (s *service) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *service) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	return s.repo.GetRoleByID(ctx, roleID)
}

func (s *service) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *service) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.repo.AssignRole(ctx, userID, roleID)
}

func (s *service) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.repo.RemoveRole(ctx, userID, roleID)
}

func (s *service) LoadUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	return s.repo.LoadUserRoles(ctx, userID)
}
