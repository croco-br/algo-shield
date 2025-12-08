package rules

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
)

func TestEvaluateAmountRule(t *testing.T) {
	engine := &Engine{}

	tests := []struct {
		name     string
		event    models.TransactionEvent
		rule     models.Rule
		expected bool
	}{
		{
			name: "Amount exceeds threshold",
			event: models.TransactionEvent{
				ExternalID:  "txn-001",
				Amount:      10000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "transfer",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "High Amount Rule",
				Type:     models.RuleTypeAmount,
				Action:   models.ActionReview,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"amount_threshold": 5000.0,
				},
				Score: 30.0,
			},
			expected: true,
		},
		{
			name: "Amount below threshold",
			event: models.TransactionEvent{
				ExternalID:  "txn-002",
				Amount:      3000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "transfer",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "High Amount Rule",
				Type:     models.RuleTypeAmount,
				Action:   models.ActionReview,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"amount_threshold": 5000.0,
				},
				Score: 30.0,
			},
			expected: false,
		},
		{
			name: "Invalid threshold type",
			event: models.TransactionEvent{
				Amount: 10000.0,
			},
			rule: models.Rule{
				Conditions: map[string]any{
					"amount_threshold": "invalid",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.evaluateAmountRule(tt.event, tt.rule)
			if result != tt.expected {
				t.Errorf("evaluateAmountRule() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateBlacklistRule(t *testing.T) {
	engine := &Engine{}

	tests := []struct {
		name     string
		event    models.TransactionEvent
		rule     models.Rule
		expected bool
	}{
		{
			name: "From account is blacklisted",
			event: models.TransactionEvent{
				ExternalID:  "txn-001",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-blacklisted",
				ToAccount:   "acc-456",
				Type:        "transfer",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "Blacklist Rule",
				Type:     models.RuleTypeBlacklist,
				Action:   models.ActionBlock,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"blacklisted_accounts": []interface{}{
						"acc-blacklisted",
						"acc-fraud",
					},
				},
				Score: 100.0,
			},
			expected: true,
		},
		{
			name: "To account is blacklisted",
			event: models.TransactionEvent{
				ExternalID:  "txn-002",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-blacklisted",
				Type:        "transfer",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "Blacklist Rule",
				Type:     models.RuleTypeBlacklist,
				Action:   models.ActionBlock,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"blacklisted_accounts": []interface{}{
						"acc-blacklisted",
						"acc-fraud",
					},
				},
				Score: 100.0,
			},
			expected: true,
		},
		{
			name: "No account is blacklisted",
			event: models.TransactionEvent{
				ExternalID:  "txn-003",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "transfer",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "Blacklist Rule",
				Type:     models.RuleTypeBlacklist,
				Action:   models.ActionBlock,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"blacklisted_accounts": []interface{}{
						"acc-blacklisted",
						"acc-fraud",
					},
				},
				Score: 100.0,
			},
			expected: false,
		},
		{
			name: "Empty blacklist",
			event: models.TransactionEvent{
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
			},
			rule: models.Rule{
				Conditions: map[string]any{
					"blacklisted_accounts": []interface{}{},
				},
			},
			expected: false,
		},
		{
			name: "Invalid blacklist format",
			event: models.TransactionEvent{
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
			},
			rule: models.Rule{
				Conditions: map[string]any{
					"blacklisted_accounts": "invalid",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.evaluateBlacklistRule(tt.event, tt.rule)
			if result != tt.expected {
				t.Errorf("evaluateBlacklistRule() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluatePatternRule(t *testing.T) {
	engine := &Engine{}

	tests := []struct {
		name     string
		event    models.TransactionEvent
		rule     models.Rule
		expected bool
	}{
		{
			name: "Pattern matches transaction type",
			event: models.TransactionEvent{
				ExternalID:  "txn-001",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "international",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "International Pattern Rule",
				Type:     models.RuleTypePattern,
				Action:   models.ActionReview,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"pattern": "international",
				},
				Score: 20.0,
			},
			expected: true,
		},
		{
			name: "Pattern does not match",
			event: models.TransactionEvent{
				ExternalID:  "txn-002",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "domestic",
			},
			rule: models.Rule{
				ID:       uuid.New(),
				Name:     "International Pattern Rule",
				Type:     models.RuleTypePattern,
				Action:   models.ActionReview,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"pattern": "international",
				},
				Score: 20.0,
			},
			expected: false,
		},
		{
			name: "Invalid pattern type",
			event: models.TransactionEvent{
				Type: "transfer",
			},
			rule: models.Rule{
				Conditions: map[string]any{
					"pattern": 123,
				},
			},
			expected: false,
		},
		{
			name: "Missing pattern",
			event: models.TransactionEvent{
				Type: "transfer",
			},
			rule: models.Rule{
				Conditions: map[string]any{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.evaluatePatternRule(tt.event, tt.rule)
			if result != tt.expected {
				t.Errorf("evaluatePatternRule() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	engine := &Engine{
		rules: []models.Rule{
			{
				ID:       uuid.New(),
				Name:     "High Amount Rule",
				Type:     models.RuleTypeAmount,
				Action:   models.ActionReview,
				Priority: 1,
				Enabled:  true,
				Conditions: map[string]any{
					"amount_threshold": 5000.0,
				},
				Score: 30.0,
			},
			{
				ID:       uuid.New(),
				Name:     "Blacklist Rule",
				Type:     models.RuleTypeBlacklist,
				Action:   models.ActionBlock,
				Priority: 2,
				Enabled:  true,
				Conditions: map[string]any{
					"blacklisted_accounts": []interface{}{
						"acc-blacklisted",
						"acc-fraud",
					},
				},
				Score: 100.0,
			},
			{
				ID:       uuid.New(),
				Name:     "International Pattern Rule",
				Type:     models.RuleTypePattern,
				Action:   models.ActionScore,
				Priority: 3,
				Enabled:  true,
				Conditions: map[string]any{
					"pattern": "international",
				},
				Score: 15.0,
			},
		},
	}

	tests := []struct {
		name              string
		event             models.TransactionEvent
		expectedStatus    models.TransactionStatus
		expectedRiskLevel models.RiskLevel
		expectedMinScore  float64
		expectedMaxScore  float64
		expectedRules     []string
	}{
		{
			name: "Low risk transaction - no rules matched",
			event: models.TransactionEvent{
				ExternalID:  "txn-001",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "domestic",
				Timestamp:   time.Now(),
			},
			expectedStatus:    models.StatusApproved,
			expectedRiskLevel: models.RiskLow,
			expectedMinScore:  0.0,
			expectedMaxScore:  0.0,
			expectedRules:     []string{},
		},
		{
			name: "Medium risk transaction - high amount",
			event: models.TransactionEvent{
				ExternalID:  "txn-002",
				Amount:      10000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "domestic",
				Timestamp:   time.Now(),
			},
			expectedStatus:    models.StatusReview,
			expectedRiskLevel: models.RiskLow,
			expectedMinScore:  30.0,
			expectedMaxScore:  30.0,
			expectedRules:     []string{"High Amount Rule"},
		},
		{
			name: "High risk transaction - blacklisted account",
			event: models.TransactionEvent{
				ExternalID:  "txn-003",
				Amount:      5000.0,
				Currency:    "USD",
				FromAccount: "acc-blacklisted",
				ToAccount:   "acc-456",
				Type:        "domestic",
				Timestamp:   time.Now(),
			},
			expectedStatus:    models.StatusRejected,
			expectedRiskLevel: models.RiskHigh,
			expectedMinScore:  100.0,
			expectedMaxScore:  100.0,
			expectedRules:     []string{"Blacklist Rule"},
		},
		{
			name: "Multiple rules matched",
			event: models.TransactionEvent{
				ExternalID:  "txn-004",
				Amount:      15000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "international",
				Timestamp:   time.Now(),
			},
			expectedStatus:    models.StatusReview,
			expectedRiskLevel: models.RiskLow,
			expectedMinScore:  45.0,
			expectedMaxScore:  45.0,
			expectedRules:     []string{"High Amount Rule", "International Pattern Rule"},
		},
		{
			name: "High risk score triggers review",
			event: models.TransactionEvent{
				ExternalID:  "txn-005",
				Amount:      20000.0,
				Currency:    "USD",
				FromAccount: "acc-blacklisted",
				ToAccount:   "acc-456",
				Type:        "international",
				Timestamp:   time.Now(),
			},
			expectedStatus:    models.StatusRejected,
			expectedRiskLevel: models.RiskHigh,
			expectedMinScore:  145.0,
			expectedMaxScore:  145.0,
			expectedRules:     []string{"High Amount Rule", "Blacklist Rule", "International Pattern Rule"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := engine.Evaluate(ctx, tt.event)

			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", result.Status, tt.expectedStatus)
			}

			if result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %v, want %v", result.RiskLevel, tt.expectedRiskLevel)
			}

			if result.RiskScore < tt.expectedMinScore || result.RiskScore > tt.expectedMaxScore {
				t.Errorf("RiskScore = %v, want between %v and %v", result.RiskScore, tt.expectedMinScore, tt.expectedMaxScore)
			}

			if len(result.MatchedRules) != len(tt.expectedRules) {
				t.Errorf("MatchedRules count = %v, want %v", len(result.MatchedRules), len(tt.expectedRules))
			}

			// Verify matched rules
			for i, ruleName := range tt.expectedRules {
				if i >= len(result.MatchedRules) {
					t.Errorf("Expected rule %s not found in matched rules", ruleName)
					continue
				}
				if result.MatchedRules[i] != ruleName {
					t.Errorf("MatchedRules[%d] = %v, want %v", i, result.MatchedRules[i], ruleName)
				}
			}

			if result.ProcessingTime < 0 {
				t.Errorf("ProcessingTime = %v, should be non-negative", result.ProcessingTime)
			}
		})
	}
}

func TestEvaluateRiskLevelThresholds(t *testing.T) {
	tests := []struct {
		name              string
		score             float64
		expectedRiskLevel models.RiskLevel
		initialStatus     models.TransactionStatus
		expectedStatus    models.TransactionStatus
	}{
		{
			name:              "Low risk - score below 50",
			score:             30.0,
			expectedRiskLevel: models.RiskLow,
			initialStatus:     models.StatusApproved,
			expectedStatus:    models.StatusApproved,
		},
		{
			name:              "Medium risk - score 50-79",
			score:             60.0,
			expectedRiskLevel: models.RiskMedium,
			initialStatus:     models.StatusApproved,
			expectedStatus:    models.StatusApproved,
		},
		{
			name:              "High risk - score 80+",
			score:             85.0,
			expectedRiskLevel: models.RiskHigh,
			initialStatus:     models.StatusApproved,
			expectedStatus:    models.StatusReview,
		},
		{
			name:              "High risk - exact threshold",
			score:             80.0,
			expectedRiskLevel: models.RiskHigh,
			initialStatus:     models.StatusApproved,
			expectedStatus:    models.StatusReview,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &Engine{
				rules: []models.Rule{
					{
						ID:       uuid.New(),
						Name:     "Test Rule",
						Type:     models.RuleTypeAmount,
						Action:   models.ActionScore,
						Priority: 1,
						Enabled:  true,
						Conditions: map[string]any{
							"amount_threshold": 0.0, // Always matches
						},
						Score: tt.score,
					},
				},
			}

			ctx := context.Background()
			event := models.TransactionEvent{
				ExternalID:  "txn-test",
				Amount:      1000.0,
				Currency:    "USD",
				FromAccount: "acc-123",
				ToAccount:   "acc-456",
				Type:        "transfer",
				Timestamp:   time.Now(),
			}

			result, err := engine.Evaluate(ctx, event)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}

			if result.RiskLevel != tt.expectedRiskLevel {
				t.Errorf("RiskLevel = %v, want %v", result.RiskLevel, tt.expectedRiskLevel)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", result.Status, tt.expectedStatus)
			}
		})
	}
}

func TestNewEngine(t *testing.T) {
	engine := NewEngine(nil, nil)

	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	if engine.rules == nil {
		t.Error("NewEngine() rules should be initialized")
	}

	if len(engine.rules) != 0 {
		t.Errorf("NewEngine() rules length = %v, want 0", len(engine.rules))
	}
}
