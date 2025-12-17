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
	LoadUserRoles(ctx context.Context, user *models.User) error
	LoadUserGroups(ctx context.Context, user *models.User) error
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

	// Load roles and groups
	if err := r.LoadUserRoles(ctx, &user); err != nil {
		return nil, err
	}
	if err := r.LoadUserGroups(ctx, &user); err != nil {
		return nil, err
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

	// Load roles and groups
	if err := r.LoadUserRoles(ctx, &user); err != nil {
		return nil, err
	}
	if err := r.LoadUserGroups(ctx, &user); err != nil {
		return nil, err
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

func (r *PostgresUserRepository) LoadUserRoles(ctx context.Context, user *models.User) error {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		UNION
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN group_roles gr ON r.id = gr.role_id
		INNER JOIN user_groups ug ON gr.group_id = ug.group_id
		WHERE ug.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	roleMap := make(map[uuid.UUID]bool) // To avoid duplicates

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			continue
		}
		if !roleMap[role.ID] {
			roles = append(roles, role)
			roleMap[role.ID] = true
		}
	}

	user.Roles = roles
	return nil
}

func (r *PostgresUserRepository) LoadUserGroups(ctx context.Context, user *models.User) error {
	query := `
		SELECT g.id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt); err != nil {
			continue
		}
		groups = append(groups, group)
	}

	user.Groups = groups
	return nil
}

// RoleRepository defines the interface for role operations
type RoleRepository interface {
	GetRoleIDByName(ctx context.Context, tx pgx.Tx, name string) (uuid.UUID, error)
	AssignRoleToUser(ctx context.Context, tx pgx.Tx, userID, roleID uuid.UUID) error
}

// PostgresRoleRepository is the PostgreSQL implementation of RoleRepository
type PostgresRoleRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &PostgresRoleRepository{db: db}
}

func (r *PostgresRoleRepository) GetRoleIDByName(ctx context.Context, tx pgx.Tx, name string) (uuid.UUID, error) {
	var roleID uuid.UUID
	query := `SELECT id FROM roles WHERE name = $1 LIMIT 1`
	err := tx.QueryRow(ctx, query, name).Scan(&roleID)
	return roleID, err
}

func (r *PostgresRoleRepository) AssignRoleToUser(ctx context.Context, tx pgx.Tx, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := tx.Exec(ctx, query, userID, roleID, time.Now())
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
