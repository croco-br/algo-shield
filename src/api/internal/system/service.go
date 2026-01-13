package system

import (
	"context"
)

// Service defines the interface for system business logic
type Service interface {
	GetSyntheticMode(ctx context.Context) (*SyntheticModeResponse, error)
	SetSyntheticMode(ctx context.Context, enabled bool) (*SyntheticModeResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new system service with dependency injection
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetSyntheticMode retrieves the current synthetic mode status
func (s *service) GetSyntheticMode(ctx context.Context) (*SyntheticModeResponse, error) {
	config, updatedAt, err := s.repo.GetSyntheticMode(ctx)
	if err != nil {
		return nil, err
	}

	return &SyntheticModeResponse{
		Enabled:   config.Enabled,
		UpdatedAt: updatedAt,
	}, nil
}

// SetSyntheticMode updates the synthetic mode status
func (s *service) SetSyntheticMode(ctx context.Context, enabled bool) (*SyntheticModeResponse, error) {
	if err := s.repo.SetSyntheticMode(ctx, enabled); err != nil {
		return nil, err
	}

	return s.GetSyntheticMode(ctx)
}
