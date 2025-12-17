package rules

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Repository defines the interface for rule data access operations
type Repository interface {
	// LoadRules loads enabled rules from database or cache (used by worker)
	LoadRules(ctx context.Context) ([]models.Rule, error)
	// CreateRule creates a new rule
	CreateRule(ctx context.Context, rule *models.Rule) error
	// GetRule retrieves a rule by ID
	GetRule(ctx context.Context, id uuid.UUID) (*models.Rule, error)
	// ListRules retrieves all rules
	ListRules(ctx context.Context) ([]models.Rule, error)
	// UpdateRule updates an existing rule
	UpdateRule(ctx context.Context, rule *models.Rule) error
	// DeleteRule deletes a rule by ID
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewPostgresRepository creates a new PostgreSQL rule repository
func NewPostgresRepository(db *pgxpool.Pool, redis *redis.Client) Repository {
	return &PostgresRepository{db: db, redis: redis}
}

// LoadRules loads enabled rules from database or cache (used by worker)
func (r *PostgresRepository) LoadRules(ctx context.Context) ([]models.Rule, error) {
	// Try to get from cache first
	if r.redis != nil {
		cachedRules, err := r.redis.Get(ctx, "rules:cache").Result()
		if err == nil && cachedRules != "" {
			var rules []models.Rule
			if err := json.Unmarshal([]byte(cachedRules), &rules); err == nil {
				return rules, nil
			}
		}
	}

	// Load from database
	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		WHERE enabled = true
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
			log.Printf("Failed to unmarshal rule conditions for rule %s: %v", rule.Name, err)
			continue
		}
		rules = append(rules, rule)
	}

	// Cache rules if redis is available
	if r.redis != nil {
		rulesJSON, err := json.Marshal(rules)
		if err != nil {
			log.Printf("Failed to marshal rules for cache: %v", err)
		} else {
			r.redis.Set(ctx, "rules:cache", rulesJSON, 5*time.Minute)
		}
	}

	return rules, nil
}

// CreateRule creates a new rule
func (r *PostgresRepository) CreateRule(ctx context.Context, rule *models.Rule) error {
	conditionsJSON, err := json.Marshal(rule.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

	query := `
		INSERT INTO rules (id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = r.db.Exec(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Type, rule.Action,
		rule.Priority, rule.Enabled, conditionsJSON, rule.Score,
		rule.CreatedAt, rule.UpdatedAt,
	)

	// Invalidate cache on create
	if err == nil && r.redis != nil {
		r.redis.Del(ctx, "rules:cache")
	}

	return err
}

// GetRule retrieves a rule by ID
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	if err := json.Unmarshal(conditionsJSON, &rule.Conditions); err != nil {
		return nil, err
	}

	return &rule, nil
}

// ListRules retrieves all rules
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

// UpdateRule updates an existing rule
func (r *PostgresRepository) UpdateRule(ctx context.Context, rule *models.Rule) error {
	conditionsJSON, err := json.Marshal(rule.Conditions)
	if err != nil {
		return fmt.Errorf("failed to marshal conditions: %w", err)
	}

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
		return pgx.ErrNoRows
	}

	// Invalidate cache on update
	if r.redis != nil {
		r.redis.Del(ctx, "rules:cache")
	}

	return nil
}

// DeleteRule deletes a rule by ID
func (r *PostgresRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM rules WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	// Invalidate cache on delete
	if r.redis != nil {
		r.redis.Del(ctx, "rules:cache")
	}

	return nil
}
