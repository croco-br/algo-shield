package engine

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Repository defines the interface for rule data access operations
type Repository interface {
	// LoadRules loads rules from database or cache
	LoadRules(ctx context.Context) ([]models.Rule, error)
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

func (r *PostgresRepository) LoadRules(ctx context.Context) ([]models.Rule, error) {
	// Try to get from cache first
	cachedRules, err := r.redis.Get(ctx, "rules:cache").Result()
	if err == nil && cachedRules != "" {
		var rules []models.Rule
		if err := json.Unmarshal([]byte(cachedRules), &rules); err == nil {
			return rules, nil
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

	// Cache rules
	rulesJSON, _ := json.Marshal(rules)
	r.redis.Set(ctx, "rules:cache", rulesJSON, 5*time.Minute)

	return rules, nil
}
