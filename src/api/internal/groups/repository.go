package groups

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for group operations
type Repository interface {
	ListGroups(ctx context.Context) ([]models.Group, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetGroupByName(ctx context.Context, name string) (*models.Group, error)
	LoadUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error)
	LoadGroupRoles(ctx context.Context, groupID uuid.UUID) ([]models.Role, error)
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListGroups(ctx context.Context) ([]models.Group, error) {
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

func (r *PostgresRepository) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, groupID).Scan(
		&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *PostgresRepository) GetGroupByName(ctx context.Context, name string) (*models.Group, error) {
	var group models.Group
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		WHERE name = $1
	`

	err := r.db.QueryRow(ctx, query, name).Scan(
		&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *PostgresRepository) LoadUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
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

func (r *PostgresRepository) LoadGroupRoles(ctx context.Context, groupID uuid.UUID) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN group_roles gr ON r.id = gr.role_id
		WHERE gr.group_id = $1
	`

	rows, err := r.db.Query(ctx, query, groupID)
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
