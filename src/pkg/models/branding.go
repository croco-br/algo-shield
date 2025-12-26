package models

import "time"

// BrandingConfig represents the white label branding configuration
type BrandingConfig struct {
	ID             int       `json:"id"`
	AppName        string    `json:"app_name"`
	IconURL        *string   `json:"icon_url,omitempty"`
	FaviconURL     *string   `json:"favicon_url,omitempty"`
	PrimaryColor   string    `json:"primary_color"`
	SecondaryColor string    `json:"secondary_color"`
	HeaderColor    string    `json:"header_color"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
