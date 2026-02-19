package permissions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
)

// TokenRevoker revokes all tokens for a user (e.g., on deactivation).
// Implemented by auth.Service to avoid circular imports.
type TokenRevoker interface {
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
}

// PermissionsService defines the interface for permissions operations
type PermissionsService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUserActive(ctx context.Context, currentUserID, targetUserID uuid.UUID, active bool) error
}

// Service manages users with their roles and groups aggregated.
// This service has a clear responsibility: user management with permission context.
type Service struct {
	userRepo     UserRepository
	roleService  roles.Service  // Only used for LoadUserRoles
	groupService groups.Service // Only used for LoadUserGroups
	tokenRevoker TokenRevoker   // Optional: revokes tokens on user deactivation
}

// NewService creates a new permissions service with dependency injection
// Follows Dependency Inversion Principle - receives interfaces, not concrete types
// tokenRevoker can be nil; if set, tokens are revoked when a user is deactivated
func NewService(userRepo UserRepository, roleService roles.Service, groupService groups.Service, tokenRevoker TokenRevoker) *Service {
	return &Service{
		userRepo:     userRepo,
		roleService:  roleService,
		groupService: groupService,
		tokenRevoker: tokenRevoker,
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

// UpdateUserActive updates a user's active status with admin protection
// Prevents admin from deactivating themselves or the last active admin
func (s *Service) UpdateUserActive(ctx context.Context, currentUserID, targetUserID uuid.UUID, active bool) error {
	// If activating, no need for protection checks
	if active {
		return s.userRepo.UpdateUserActive(ctx, targetUserID, active)
	}

	// Protection 1: Prevent admin from deactivating themselves
	if currentUserID == targetUserID {
		return apierrors.CannotDeactivateSelf()
	}

	// Protection 2: Check if target user is an admin
	isTargetAdmin, err := s.userRepo.HasAdminRole(ctx, targetUserID)
	if err != nil {
		return apierrors.InternalError("Failed to check user role")
	}

	// If target is admin, check if they are the last active admin
	if isTargetAdmin {
		activeAdminCount, err := s.userRepo.CountActiveAdmins(ctx, &targetUserID)
		if err != nil {
			return apierrors.InternalError("Failed to count active administrators")
		}

		if activeAdminCount == 0 {
			return apierrors.CannotDeactivateLastAdmin()
		}
	}

	// Validation passed, update user status
	if err := s.userRepo.UpdateUserActive(ctx, targetUserID, active); err != nil {
		return err
	}

	// SECURITY: Revoke all tokens for the deactivated user so they cannot continue using the app
	if s.tokenRevoker != nil {
		if err := s.tokenRevoker.RevokeAllUserTokens(ctx, targetUserID); err != nil {
			return apierrors.InternalError("Failed to revoke user tokens")
		}
	}
	return nil
}

// loadUserRolesAndGroups aggregates roles and groups for a user.
// This is the core responsibility of this service: enriching user data with permission context.
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
