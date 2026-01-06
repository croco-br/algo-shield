package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	StatusPending  TransactionStatus = "pending"
	StatusApproved TransactionStatus = "approved"
	StatusRejected TransactionStatus = "rejected"
	StatusInReview TransactionStatus = "in_review"
)

type Transaction struct {
	ID             uuid.UUID         `json:"id"`
	ExternalID     string            `json:"external_id"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	Origin         string            `json:"origin"`
	Destination    string            `json:"destination"`
	Type           string            `json:"type"`
	Status         TransactionStatus `json:"status"`
	ProcessingTime int64             `json:"processing_time_ms"`
	MatchedRules   []string          `json:"matched_rules"`
	Metadata       map[string]any    `json:"metadata"`
	CreatedAt      time.Time         `json:"created_at"`
	ProcessedAt    *time.Time        `json:"processed_at"`
}

// Event represents a generic JSON event for rule evaluation
// The structure is defined by the event schema, not hardcoded
type Event map[string]any

type TransactionResult struct {
	TransactionID  uuid.UUID         `json:"transaction_id"`
	Status         TransactionStatus `json:"status"`
	MatchedRules   []string          `json:"matched_rules"`
	ProcessingTime int64             `json:"processing_time_ms"`
	Message        string            `json:"message"`
}
