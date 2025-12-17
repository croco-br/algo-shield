package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	userRepo  UserRepository
	roleRepo  RoleRepository
	txManager TransactionManager
}

func NewService(db *pgxpool.Pool, cfg *config.Config) *Service {
	return &Service{
		userRepo:  NewPostgresUserRepository(db),
		roleRepo:  NewPostgresRoleRepository(db),
		txManager: NewPostgresTransactionManager(db),
	}
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email, false)
}

func (s *Service) GetUserByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email, true)
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
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
	if err := s.userRepo.LoadUserRoles(ctx, user); err != nil {
		return nil, err
	}
	if err := s.userRepo.LoadUserGroups(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error {
	return s.userRepo.UpdateLastLogin(ctx, userID, lastLoginAt)
}
