package permissions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	userRepo  UserRepository
	roleRepo  RoleRepository
	groupRepo GroupRepository
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		userRepo:  NewPostgresUserRepository(db),
		roleRepo:  NewPostgresRoleRepository(db),
		groupRepo: NewPostgresGroupRepository(db),
	}
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *Service) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.ListUsers(ctx)
}

func (s *Service) UpdateUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	return s.userRepo.UpdateUserActive(ctx, userID, active)
}

func (s *Service) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.AssignRole(ctx, userID, roleID)
}

func (s *Service) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.RemoveRole(ctx, userID, roleID)
}

func (s *Service) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.roleRepo.ListRoles(ctx)
}

func (s *Service) ListGroups(ctx context.Context) ([]models.Group, error) {
	return s.groupRepo.ListGroups(ctx)
}
