package branding

import (
	"context"
	"encoding/json"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Repository defines the interface for branding configuration operations
type Repository interface {
	Get(ctx context.Context) (*models.BrandingConfig, error)
	Update(ctx context.Context, config *models.BrandingConfig) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewPostgresRepository creates a new PostgreSQL repository for branding operations
func NewPostgresRepository(db *pgxpool.Pool, redis *redis.Client) Repository {
	return &PostgresRepository{db: db, redis: redis}
}

// Get retrieves the current branding configuration
// Uses Redis cache to minimize database queries
func (r *PostgresRepository) Get(ctx context.Context) (*models.BrandingConfig, error) {
	// Try to get from cache first
	if r.redis != nil {
		cachedConfig, err := r.redis.Get(ctx, "branding:config").Result()
		if err == nil && cachedConfig != "" {
			var config models.BrandingConfig
			if err := json.Unmarshal([]byte(cachedConfig), &config); err == nil {
				return &config, nil
			}
		}
	}

	// Load from database
	var config models.BrandingConfig
	query := `
		SELECT id, app_name, icon_url, favicon_url, primary_color, secondary_color, header_color, created_at, updated_at
		FROM branding_config
		WHERE id = 1
	`

	err := r.db.QueryRow(ctx, query).Scan(
		&config.ID,
		&config.AppName,
		&config.IconURL,
		&config.FaviconURL,
		&config.PrimaryColor,
		&config.SecondaryColor,
		&config.HeaderColor,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Cache the result if redis is available
	if r.redis != nil {
		configJSON, err := json.Marshal(config)
		if err == nil {
			// Cache for 10 minutes (branding changes infrequently)
			r.redis.Set(ctx, "branding:config", configJSON, 10*time.Minute)
		}
	}

	return &config, nil
}

// Update updates the branding configuration
// Invalidates cache after successful update
func (r *PostgresRepository) Update(ctx context.Context, config *models.BrandingConfig) error {
	query := `
		UPDATE branding_config
		SET app_name = $1,
		    icon_url = $2,
		    favicon_url = $3,
		    primary_color = $4,
		    secondary_color = $5,
		    header_color = $6,
		    updated_at = $7
		WHERE id = 1
		RETURNING id, app_name, icon_url, favicon_url, primary_color, secondary_color, header_color, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		config.AppName,
		config.IconURL,
		config.FaviconURL,
		config.PrimaryColor,
		config.SecondaryColor,
		config.HeaderColor,
		now,
	).Scan(
		&config.ID,
		&config.AppName,
		&config.IconURL,
		&config.FaviconURL,
		&config.PrimaryColor,
		&config.SecondaryColor,
		&config.HeaderColor,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	// Invalidate cache on successful update
	if err == nil && r.redis != nil {
		r.redis.Del(ctx, "branding:config")
	}

	return err
}
