package rules

import (
	"context"
	"encoding/json"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for rule data access operations
type Repository interface {
	CreateRule(ctx context.Context, rule *models.Rule) error
	GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error)
	ListRules(ctx context.Context) ([]models.Rule, error)
	UpdateRule(ctx context.Context, rule *models.Rule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateRule(ctx context.Context, rule *models.Rule) error {
	conditionsJSON, _ := json.Marshal(rule.Conditions)

	query := `
		INSERT INTO rules (id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Type, rule.Action,
		rule.Priority, rule.Enabled, conditionsJSON, rule.Score,
		rule.CreatedAt, rule.UpdatedAt,
	)

	return err
}

func (r *PostgresRepository) GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error) {
	var rule models.Rule
	var conditionsJSON []byte

	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Action,
		&rule.Priority, &rule.Enabled, &conditionsJSON, &rule.Score,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r *PostgresRepository) ListRules(ctx context.Context) ([]models.Rule, error) {
	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		ORDER BY priority ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]models.Rule, 0)
	for rows.Next() {
		var rule models.Rule
		var conditionsJSON []byte

		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Action,
			&rule.Priority, &rule.Enabled, &conditionsJSON, &rule.Score,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
			continue
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *PostgresRepository) UpdateRule(ctx context.Context, rule *models.Rule) error {
	conditionsJSON, _ := json.Marshal(rule.Conditions)

	query := `
		UPDATE rules
		SET name = $2, description = $3, type = $4, action = $5, 
		    priority = $6, enabled = $7, conditions = $8, score = $9, updated_at = $10
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Type, rule.Action,
		rule.Priority, rule.Enabled, conditionsJSON, rule.Score, rule.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return nil
}

func (r *PostgresRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rules WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return nil
}
