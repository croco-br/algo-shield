package engine

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// RuleLoader handles rule loading and caching
type RuleLoader struct {
	repo  Repository
	rules []models.Rule
}

// NewRuleLoader creates a new rule loader
func NewRuleLoader(db *pgxpool.Pool, redis *redis.Client) *RuleLoader {
	return &RuleLoader{
		repo:  NewPostgresRepository(db, redis),
		rules: make([]models.Rule, 0),
	}
}

// LoadRules loads rules from the repository
func (rl *RuleLoader) LoadRules(ctx context.Context) error {
	rules, err := rl.repo.LoadRules(ctx)
	if err != nil {
		return err
	}
	rl.rules = rules
	return nil
}

// GetRules returns the currently loaded rules
func (rl *RuleLoader) GetRules() []models.Rule {
	return rl.rules
}
