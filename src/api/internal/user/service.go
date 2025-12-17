package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	userRepo     UserRepository
	roleRepo     roles.Repository
	roleService  roles.Service
	groupService groups.Service
	txManager    TransactionManager
}

func NewService(db *pgxpool.Pool, cfg *config.Config, roleService roles.Service, groupService groups.Service) *Service {
	return &Service{
		userRepo:     NewPostgresUserRepository(db),
		roleRepo:     roles.NewPostgresRepository(db),
		roleService:  roleService,
		groupService: groupService,
		txManager:    NewPostgresTransactionManager(db),
	}
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email, false)
	if err != nil {
		return nil, err
	}
	return s.loadUserRolesAndGroups(ctx, user)
}

func (s *Service) GetUserByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email, true)
	if err != nil {
		return nil, err
	}
	return s.loadUserRolesAndGroups(ctx, user)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.loadUserRolesAndGroups(ctx, user)
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

func (s *Service) CreateUser(ctx context.Context, email, name, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: &passwordHash,
		AuthType:     models.AuthTypeLocal,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Start transaction
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != sql.ErrTxDone {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Insert user using repository
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Assign default viewer role to new users (by name)
	// This is non-critical - user creation succeeds even if role assignment fails
	viewerRoleID, err := s.roleRepo.GetRoleIDByName(ctx, tx, "viewer")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Warning: 'viewer' role not found, user %s created without default role", user.ID)
		} else {
			log.Printf("Error querying for 'viewer' role for user %s: %v", user.ID, err)
		}
	} else {
		// Role found, assign it within the transaction
		if err := s.roleRepo.AssignRoleToUser(ctx, tx, user.ID, viewerRoleID); err != nil {
			log.Printf("Error assigning default 'viewer' role to user %s: %v", user.ID, err)
			// Continue with user creation even if role assignment fails
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Load roles and groups (read-only operations, outside transaction)
	return s.loadUserRolesAndGroups(ctx, user)
}

func (s *Service) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error {
	return s.userRepo.UpdateLastLogin(ctx, userID, lastLoginAt)
}
