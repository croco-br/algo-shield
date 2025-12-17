package rules

import (
	"context"
	"sync"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// RuleService handles rule loading and caching for the worker
type RuleService struct {
	mu    sync.RWMutex
	repo  rules.Repository
	rules []models.Rule
}

// NewRuleService creates a new rule service for the worker
func NewRuleService(db *pgxpool.Pool, redis *redis.Client) *RuleService {
	return &RuleService{
		repo:  rules.NewPostgresRepository(db, redis),
		rules: make([]models.Rule, 0),
	}
}

// LoadRules loads rules from the repository
func (rl *RuleService) LoadRules(ctx context.Context) error {
	rules, err := rl.repo.LoadRules(ctx)
	if err != nil {
		return err
	}
	rl.mu.Lock()
	rl.rules = rules
	rl.mu.Unlock()
	return nil
}

// GetRules returns the currently loaded rules (returns a copy to avoid race conditions)
func (rl *RuleService) GetRules() []models.Rule {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	// Return a copy to prevent external modification
	return append([]models.Rule(nil), rl.rules...)
}
