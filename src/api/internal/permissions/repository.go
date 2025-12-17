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

		// Load roles and groups
		if err := r.LoadUserRoles(ctx, &user); err != nil {
			return nil, err
		}
		if err := r.LoadUserGroups(ctx, &user); err != nil {
			return nil, err
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
	ListRoles(ctx context.Context) ([]models.Role, error)
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error
}

// PostgresRoleRepository is the PostgreSQL implementation of RoleRepository
type PostgresRoleRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &PostgresRoleRepository{db: db}
}

func (r *PostgresRoleRepository) ListRoles(ctx context.Context) ([]models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *PostgresRoleRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, roleID, time.Now())
	return err
}

func (r *PostgresRoleRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

// GroupRepository defines the interface for group operations
type GroupRepository interface {
	ListGroups(ctx context.Context) ([]models.Group, error)
}

// PostgresGroupRepository is the PostgreSQL implementation of GroupRepository
type PostgresGroupRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGroupRepository(db *pgxpool.Pool) GroupRepository {
	return &PostgresGroupRepository{db: db}
}

func (r *PostgresGroupRepository) ListGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
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

	return groups, nil
}
