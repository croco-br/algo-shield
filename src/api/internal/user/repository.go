package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string, includePassword bool) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error
}

// PostgresUserRepository is the PostgreSQL implementation of UserRepository
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string, includePassword bool) (*models.User, error) {
	var user models.User
	var passwordHash, googleID, pictureURL sql.NullString
	var lastLoginAt sql.NullTime

	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
		&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if includePassword && passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if pictureURL.Valid {
		user.PictureURL = &pictureURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	var passwordHash, googleID, pictureURL sql.NullString
	var lastLoginAt sql.NullTime

	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
		&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if pictureURL.Valid {
		user.PictureURL = &pictureURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash,
		user.AuthType, user.Active, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, lastLoginAt, time.Now(), userID)
	return err
}

// TransactionManager defines the interface for database transaction operations
type TransactionManager interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// PostgresTransactionManager is the PostgreSQL implementation of TransactionManager
type PostgresTransactionManager struct {
	db interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
}

func NewPostgresTransactionManager(db interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}) TransactionManager {
	return &PostgresTransactionManager{db: db}
}

func (m *PostgresTransactionManager) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.db.Begin(ctx)
}
