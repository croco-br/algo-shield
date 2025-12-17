package permissions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	userRepo     UserRepository
	roleService  roles.Service
	groupService groups.Service
}

func NewService(db *pgxpool.Pool, roleService roles.Service, groupService groups.Service) *Service {
	return &Service{
		userRepo:     NewPostgresUserRepository(db),
		roleService:  roleService,
		groupService: groupService,
	}
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.loadUserRolesAndGroups(ctx, user)
}

func (s *Service) ListUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Load roles and groups for each user
	for i := range users {
		user, err := s.loadUserRolesAndGroups(ctx, &users[i])
		if err != nil {
			return nil, err
		}
		users[i] = *user
	}

	return users, nil
}

func (s *Service) UpdateUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	return s.userRepo.UpdateUserActive(ctx, userID, active)
}

func (s *Service) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleService.AssignRole(ctx, userID, roleID)
}

func (s *Service) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleService.RemoveRole(ctx, userID, roleID)
}

func (s *Service) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.roleService.ListRoles(ctx)
}

func (s *Service) ListGroups(ctx context.Context) ([]models.Group, error) {
	return s.groupService.ListGroups(ctx)
}

func (s *Service) loadUserRolesAndGroups(ctx context.Context, user *models.User) (*models.User, error) {
	// Load roles
	roles, err := s.roleService.LoadUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	// Load groups
	groups, err := s.groupService.LoadUserGroups(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Groups = groups

	return user, nil
}
