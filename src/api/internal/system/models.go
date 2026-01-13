package system

import "time"

// SystemConfig represents a system configuration entry
type SystemConfig struct {
	Key       string         `json:"key"`
	Value     map[string]any `json:"value"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// SyntheticModeConfig represents the synthetic mode configuration
type SyntheticModeConfig struct {
	Enabled bool `json:"enabled"`
}

// SyntheticModeResponse is the API response for synthetic mode status
type SyntheticModeResponse struct {
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateSyntheticModeRequest is the request to update synthetic mode
type UpdateSyntheticModeRequest struct {
	Enabled bool `json:"enabled"`
}
