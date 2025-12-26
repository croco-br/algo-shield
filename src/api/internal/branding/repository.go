package branding

import (
	"context"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for branding configuration operations
type Repository interface {
	Get(ctx context.Context) (*models.BrandingConfig, error)
	Update(ctx context.Context, config *models.BrandingConfig) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository for branding operations
func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

// Get retrieves the current branding configuration
func (r *PostgresRepository) Get(ctx context.Context) (*models.BrandingConfig, error) {
	var config models.BrandingConfig
	query := `
		SELECT id, app_name, icon_url, favicon_url, primary_color, secondary_color, created_at, updated_at
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
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Update updates the branding configuration
func (r *PostgresRepository) Update(ctx context.Context, config *models.BrandingConfig) error {
	query := `
		UPDATE branding_config
		SET app_name = $1,
		    icon_url = $2,
		    favicon_url = $3,
		    primary_color = $4,
		    secondary_color = $5,
		    updated_at = $6
		WHERE id = 1
		RETURNING id, app_name, icon_url, favicon_url, primary_color, secondary_color, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		config.AppName,
		config.IconURL,
		config.FaviconURL,
		config.PrimaryColor,
		config.SecondaryColor,
		now,
	).Scan(
		&config.ID,
		&config.AppName,
		&config.IconURL,
		&config.FaviconURL,
		&config.PrimaryColor,
		&config.SecondaryColor,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	return err
}
