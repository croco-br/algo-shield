package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthType string

const (
	AuthTypeSSO   AuthType = "sso"
	AuthTypeLocal AuthType = "local"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	PasswordHash *string    `json:"-"` // Never send password hash in JSON
	GoogleID     *string    `json:"google_id,omitempty"`
	PictureURL   *string    `json:"picture_url,omitempty"`
	AuthType     AuthType   `json:"auth_type"`
	Active       bool       `json:"active"`
	Roles        []Role     `json:"roles,omitempty"`
	Groups       []Group    `json:"groups,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Roles       []Role    `json:"roles,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
