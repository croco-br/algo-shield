package system

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	syntheticModeKey      = "synthetic_mode"
	syntheticModeCacheKey = "system:synthetic_mode"
	cacheTTL              = 5 * time.Minute
)

// Repository defines the interface for system configuration data access
type Repository interface {
	GetSyntheticMode(ctx context.Context) (*SyntheticModeConfig, time.Time, error)
	SetSyntheticMode(ctx context.Context, enabled bool) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewPostgresRepository creates a new PostgreSQL system repository
func NewPostgresRepository(db *pgxpool.Pool, redis *redis.Client) Repository {
	return &PostgresRepository{
		db:    db,
		redis: redis,
	}
}

// GetSyntheticMode retrieves the synthetic mode configuration
func (r *PostgresRepository) GetSyntheticMode(ctx context.Context) (*SyntheticModeConfig, time.Time, error) {
	// Try cache first
	cached, updatedAt, err := r.getFromCache(ctx)
	if err == nil && cached != nil {
		return cached, updatedAt, nil
	}

	// Query database
	var valueJSON []byte
	var updatedAtDB time.Time
	query := `SELECT value, updated_at FROM system_config WHERE key = $1`
	err = r.db.QueryRow(ctx, query, syntheticModeKey).Scan(&valueJSON, &updatedAtDB)
	if err != nil {
		// Return default if not found
		return &SyntheticModeConfig{Enabled: false}, time.Now(), nil
	}

	var config SyntheticModeConfig
	if err := json.Unmarshal(valueJSON, &config); err != nil {
		return &SyntheticModeConfig{Enabled: false}, updatedAtDB, nil
	}

	// Update cache
	_ = r.setCache(ctx, &config, updatedAtDB)

	return &config, updatedAtDB, nil
}

// SetSyntheticMode updates the synthetic mode configuration
func (r *PostgresRepository) SetSyntheticMode(ctx context.Context, enabled bool) error {
	config := SyntheticModeConfig{Enabled: enabled}
	valueJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	now := time.Now()
	query := `
		INSERT INTO system_config (key, value, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = $3
	`
	_, err = r.db.Exec(ctx, query, syntheticModeKey, valueJSON, now)
	if err != nil {
		return err
	}

	// Update cache
	_ = r.setCache(ctx, &config, now)

	return nil
}

func (r *PostgresRepository) getFromCache(ctx context.Context) (*SyntheticModeConfig, time.Time, error) {
	data, err := r.redis.Get(ctx, syntheticModeCacheKey).Bytes()
	if err != nil {
		return nil, time.Time{}, err
	}

	var cached struct {
		Config    SyntheticModeConfig `json:"config"`
		UpdatedAt time.Time           `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, time.Time{}, err
	}

	return &cached.Config, cached.UpdatedAt, nil
}

func (r *PostgresRepository) setCache(ctx context.Context, config *SyntheticModeConfig, updatedAt time.Time) error {
	cached := struct {
		Config    SyntheticModeConfig `json:"config"`
		UpdatedAt time.Time           `json:"updated_at"`
	}{
		Config:    *config,
		UpdatedAt: updatedAt,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, syntheticModeCacheKey, data, cacheTTL).Err()
}
