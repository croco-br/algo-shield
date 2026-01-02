package models

import (
	"time"

	"github.com/google/uuid"
)

type RuleType string

const (
	RuleTypeAmount    RuleType = "amount"
	RuleTypeVelocity  RuleType = "velocity"
	RuleTypeBlocklist RuleType = "blocklist"
	RuleTypeGeography RuleType = "geography"
	RuleTypeCustom    RuleType = "custom"
)

type RuleAction string

const (
	ActionAllow  RuleAction = "allow"
	ActionBlock  RuleAction = "block"
	ActionReview RuleAction = "review"
	ActionScore  RuleAction = "score"
)

type Rule struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name" validate:"required,min=1,max=255"`
	Description string         `json:"description" validate:"max=1000"`
	Type        RuleType       `json:"type" validate:"required,oneof=amount velocity blocklist geography custom"`
	Action      RuleAction     `json:"action" validate:"required,oneof=allow block review score"`
	Priority    int            `json:"priority" validate:"gte=0,lte=1000"`
	Enabled     bool           `json:"enabled"`
	Conditions  map[string]any `json:"conditions" validate:"required"`
	Score       float64        `json:"score" validate:"gte=0,lte=100"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type RuleConfig struct {
	AmountThreshold     *float64 `json:"amount_threshold,omitempty"`
	TransactionCount    *int     `json:"transaction_count,omitempty"`
	TimeWindowSeconds   *int     `json:"time_window_seconds,omitempty"`
	BlocklistedAccounts []string `json:"blocklisted_accounts,omitempty"`
	CustomExpression    *string  `json:"custom_expression,omitempty"`
}
