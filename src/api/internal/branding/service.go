package branding

import (
	"context"
	"fmt"
	"regexp"

	"github.com/algo-shield/algo-shield/src/pkg/models"
)

// Default branding values
const (
	DefaultAppName        = "AlgoShield"
	DefaultPrimaryColor   = "#3B82F6"
	DefaultSecondaryColor = "#10B981"
	DefaultHeaderColor    = "#1e1e1e"
	DefaultIconURL        = "/assets/logo.svg"
	DefaultFaviconURL     = "/favicon.ico"
)

var (
	// hexColorRegex matches hex colors in format #RGB or #RRGGBB
	hexColorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
)

// Service defines the interface for branding business logic
type Service interface {
	GetBranding(ctx context.Context) (*models.BrandingConfig, error)
	UpdateBranding(ctx context.Context, req *UpdateBrandingRequest) (*models.BrandingConfig, error)
}

type service struct {
	repo Repository
}

// NewService creates a new branding service with dependency injection
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetBranding retrieves the current branding configuration
// Returns default values if no configuration exists
func (s *service) GetBranding(ctx context.Context) (*models.BrandingConfig, error) {
	config, err := s.repo.Get(ctx)
	if err != nil {
		// If no config exists, return defaults
		defaultIconURL := DefaultIconURL
		defaultFaviconURL := DefaultFaviconURL
		return &models.BrandingConfig{
			ID:             1,
			AppName:        DefaultAppName,
			IconURL:        &defaultIconURL,
			FaviconURL:     &defaultFaviconURL,
			PrimaryColor:   DefaultPrimaryColor,
			SecondaryColor: DefaultSecondaryColor,
			HeaderColor:    DefaultHeaderColor,
		}, nil
	}
	return config, nil
}

// UpdateBranding updates the branding configuration with validation
func (s *service) UpdateBranding(ctx context.Context, req *UpdateBrandingRequest) (*models.BrandingConfig, error) {
	// Validate color formats
	if err := validateHexColor(req.PrimaryColor, "primary_color"); err != nil {
		return nil, err
	}
	if err := validateHexColor(req.SecondaryColor, "secondary_color"); err != nil {
		return nil, err
	}
	if err := validateHexColor(req.HeaderColor, "header_color"); err != nil {
		return nil, err
	}

	// Validate app name length
	if len(req.AppName) > 100 {
		return nil, fmt.Errorf("app_name must be 100 characters or less")
	}
	if len(req.AppName) == 0 {
		return nil, fmt.Errorf("app_name cannot be empty")
	}

	// Create branding config from request
	config := &models.BrandingConfig{
		ID:             1,
		AppName:        req.AppName,
		IconURL:        req.IconURL,
		FaviconURL:     req.FaviconURL,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		HeaderColor:    req.HeaderColor,
	}

	// Update in database
	if err := s.repo.Update(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to update branding configuration: %w", err)
	}

	return config, nil
}

// validateHexColor validates that a color is in valid hex format (#RGB or #RRGGBB)
func validateHexColor(color, fieldName string) error {
	if !hexColorRegex.MatchString(color) {
		return fmt.Errorf("%s must be in hex format (#RGB or #RRGGBB)", fieldName)
	}
	return nil
}
