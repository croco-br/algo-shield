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
	StatusReview   TransactionStatus = "review"
)

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Transaction struct {
	ID              uuid.UUID         `json:"id"`
	ExternalID      string            `json:"external_id"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	FromAccount     string            `json:"from_account"`
	ToAccount       string            `json:"to_account"`
	Type            string            `json:"type"`
	Status          TransactionStatus `json:"status"`
	RiskScore       float64           `json:"risk_score"`
	RiskLevel       RiskLevel         `json:"risk_level"`
	ProcessingTime  int64             `json:"processing_time_ms"`
	MatchedRules    []string          `json:"matched_rules"`
	Metadata        map[string]any    `json:"metadata"`
	CreatedAt       time.Time         `json:"created_at"`
	ProcessedAt     *time.Time        `json:"processed_at"`
}

type TransactionEvent struct {
	ExternalID  string         `json:"external_id"`
	Amount      float64        `json:"amount"`
	Currency    string         `json:"currency"`
	FromAccount string         `json:"from_account"`
	ToAccount   string         `json:"to_account"`
	Type        string         `json:"type"`
	Metadata    map[string]any `json:"metadata"`
	Timestamp   time.Time      `json:"timestamp"`
}

type TransactionResult struct {
	TransactionID  uuid.UUID         `json:"transaction_id"`
	Status         TransactionStatus `json:"status"`
	RiskScore      float64           `json:"risk_score"`
	RiskLevel      RiskLevel         `json:"risk_level"`
	MatchedRules   []string          `json:"matched_rules"`
	ProcessingTime int64             `json:"processing_time_ms"`
	Message        string            `json:"message"`
}

