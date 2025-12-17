package permissions

import (
	"context"
	"database/sql"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUserActive(ctx context.Context, userID uuid.UUID, active bool) error
}

// PostgresUserRepository is the PostgreSQL implementation of UserRepository
type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{db: db}
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

func (r *PostgresUserRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		var passwordHash, googleID, pictureURL sql.NullString
		var lastLoginAt sql.NullTime

		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
			&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
		); err != nil {
			continue
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

		users = append(users, user)
	}

	return users, nil
}

func (r *PostgresUserRepository) UpdateUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	query := `UPDATE users SET active = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, active, time.Now(), userID)
	return err
}
