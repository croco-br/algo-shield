package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestTransactionStatus(t *testing.T) {
	tests := []struct {
		name   string
		status TransactionStatus
		valid  bool
	}{
		{"Pending status", StatusPending, true},
		{"Approved status", StatusApproved, true},
		{"Rejected status", StatusRejected, true},
		{"Review status", StatusReview, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.status) == 0 {
				t.Error("Status should not be empty")
			}
		})
	}
}

func TestRiskLevel(t *testing.T) {
	tests := []struct {
		name  string
		level RiskLevel
	}{
		{"Low risk", RiskLow},
		{"Medium risk", RiskMedium},
		{"High risk", RiskHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.level) == 0 {
				t.Error("Risk level should not be empty")
			}
		})
	}
}

func TestTransactionCreation(t *testing.T) {
	txn := Transaction{
		ID:         uuid.New(),
		ExternalID: "test_123",
		Amount:     100.50,
		Currency:   "USD",
		Status:     StatusPending,
		RiskLevel:  RiskLow,
	}

	if txn.ID == uuid.Nil {
		t.Error("Transaction ID should not be nil")
	}

	if txn.ExternalID != "test_123" {
		t.Errorf("Expected external_id 'test_123', got '%s'", txn.ExternalID)
	}

	if txn.Amount != 100.50 {
		t.Errorf("Expected amount 100.50, got %f", txn.Amount)
	}
}

