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
	CountActiveAdmins(ctx context.Context, excludeUserID *uuid.UUID) (int, error)
	HasAdminRole(ctx context.Context, userID uuid.UUID) (bool, error)
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

// CountActiveAdmins counts the number of active users with admin role
// excludeUserID can be used to exclude a specific user from the count
func (r *PostgresUserRepository) CountActiveAdmins(ctx context.Context, excludeUserID *uuid.UUID) (int, error) {
	var query string
	var args []interface{}

	if excludeUserID != nil {
		query = `
			SELECT COUNT(DISTINCT u.id)
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			LEFT JOIN user_groups ug ON u.id = ug.user_id
			LEFT JOIN group_roles gr ON ug.group_id = gr.group_id
			LEFT JOIN roles r2 ON gr.role_id = r2.id
			WHERE u.active = true 
			  AND u.id != $1
			  AND (r.name = 'admin' OR r2.name = 'admin')
		`
		args = []interface{}{*excludeUserID}
	} else {
		query = `
			SELECT COUNT(DISTINCT u.id)
			FROM users u
			INNER JOIN user_roles ur ON u.id = ur.user_id
			INNER JOIN roles r ON ur.role_id = r.id
			LEFT JOIN user_groups ug ON u.id = ug.user_id
			LEFT JOIN group_roles gr ON ug.group_id = gr.group_id
			LEFT JOIN roles r2 ON gr.role_id = r2.id
			WHERE u.active = true 
			  AND (r.name = 'admin' OR r2.name = 'admin')
		`
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// HasAdminRole checks if a user has the admin role (directly or through groups)
func (r *PostgresUserRepository) HasAdminRole(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_roles ur
			INNER JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.name = 'admin'
			UNION
			SELECT 1 FROM user_groups ug
			INNER JOIN group_roles gr ON ug.group_id = gr.group_id
			INNER JOIN roles r ON gr.role_id = r.id
			WHERE ug.user_id = $1 AND r.name = 'admin'
		)
	`

	var hasRole bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&hasRole)
	if err != nil {
		return false, err
	}

	return hasRole, nil
}
