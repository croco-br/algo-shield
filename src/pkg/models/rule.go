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
	RuleTypePattern   RuleType = "pattern"
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
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        RuleType       `json:"type"`
	Action      RuleAction     `json:"action"`
	Priority    int            `json:"priority"`
	Enabled     bool           `json:"enabled"`
	Conditions  map[string]any `json:"conditions"`
	Score       float64        `json:"score"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type RuleConfig struct {
	AmountThreshold     *float64 `json:"amount_threshold,omitempty"`
	TransactionCount    *int     `json:"transaction_count,omitempty"`
	TimeWindowSeconds   *int     `json:"time_window_seconds,omitempty"`
	BlocklistedAccounts []string `json:"blocklisted_accounts,omitempty"`
	Pattern             *string  `json:"pattern,omitempty"`
	CustomExpression    *string  `json:"custom_expression,omitempty"`
}
