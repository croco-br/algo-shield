package rules

import (
	"context"
	"sync/atomic"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// RuleService handles rule loading and caching for the worker
// Uses atomic.Value for lock-free reads and safe concurrent updates
type RuleService struct {
	repo  rules.Repository
	rules atomic.Value // stores []models.Rule
}

// NewRuleService creates a new rule service for the worker
func NewRuleService(db *pgxpool.Pool, redis *redis.Client) *RuleService {
	rs := &RuleService{
		repo: rules.NewPostgresRepository(db, redis),
	}
	// Initialize with empty slice
	rs.rules.Store(make([]models.Rule, 0))
	return rs
}

// LoadRules loads rules from the repository
// Thread-safe: uses atomic.Value.Store() for lock-free updates
func (rl *RuleService) LoadRules(ctx context.Context) error {
	rules, err := rl.repo.LoadRules(ctx)
	if err != nil {
		return err
	}
	// Store a copy to ensure immutability
	rulesCopy := make([]models.Rule, len(rules))
	copy(rulesCopy, rules)
	rl.rules.Store(rulesCopy)
	return nil
}

// GetRules returns the currently loaded rules
// Thread-safe: uses atomic.Value.Load() for lock-free reads
// Returns a copy to prevent external modification
func (rl *RuleService) GetRules() []models.Rule {
	rules := rl.rules.Load().([]models.Rule)
	// Return a copy to prevent external modification
	result := make([]models.Rule, len(rules))
	copy(result, rules)
	return result
}
