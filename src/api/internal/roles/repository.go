package roles

import (
	"context"
	"fmt"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for role operations
type Repository interface {
	ListRoles(ctx context.Context) ([]models.Role, error)
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	GetRoleIDByName(ctx context.Context, tx pgx.Tx, name string) (uuid.UUID, error)
	AssignRoleToUser(ctx context.Context, tx pgx.Tx, userID, roleID uuid.UUID) error
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error
	LoadUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListRoles(ctx context.Context) ([]models.Role, error) {
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

func (r *PostgresRepository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	var role models.Role
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, roleID).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *PostgresRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *PostgresRepository) GetRoleIDByName(ctx context.Context, tx pgx.Tx, name string) (uuid.UUID, error) {
	var roleID uuid.UUID
	query := `SELECT id FROM roles WHERE name = $1 LIMIT 1`
	err := tx.QueryRow(ctx, query, name).Scan(&roleID)
	if err != nil {
		return uuid.Nil, err
	}
	// Validate that we got a valid UUID
	if roleID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("role '%s' not found", name)
	}
	return roleID, nil
}

func (r *PostgresRepository) AssignRoleToUser(ctx context.Context, tx pgx.Tx, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := tx.Exec(ctx, query, userID, roleID, time.Now())
	return err
}

func (r *PostgresRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, roleID, time.Now())
	return err
}

func (r *PostgresRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *PostgresRepository) LoadUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
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

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
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

	return roles, nil
}
