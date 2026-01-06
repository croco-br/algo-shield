package models

import (
	"time"

	"github.com/google/uuid"
)

type RuleType string

const (
	RuleTypeCustom RuleType = "custom"
)

type RuleAction string

const (
	ActionAllow  RuleAction = "allow"
	ActionBlock  RuleAction = "block"
	ActionReview RuleAction = "review"
)

type Rule struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name" validate:"required,min=1,max=255"`
	Description string         `json:"description" validate:"max=1000"`
	Type        RuleType       `json:"type" validate:"required,oneof=custom"`
	Action      RuleAction     `json:"action" validate:"required,oneof=allow block review"`
	Priority    int            `json:"priority" validate:"gte=0,lte=100"`
	Enabled     bool           `json:"enabled"`
	Conditions  map[string]any `json:"conditions" validate:"required"`
	SchemaID    *uuid.UUID     `json:"schema_id,omitempty"` // Reference to event schema
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
