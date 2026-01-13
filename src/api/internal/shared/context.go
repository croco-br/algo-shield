package shared

import (
	"context"
)

type contextKey string

const (
	syntheticModeCtxKey contextKey = "synthetic_mode"
)

// WithSyntheticMode returns a new context with the synthetic mode value
func WithSyntheticMode(ctx context.Context, syntheticMode bool) context.Context {
	return context.WithValue(ctx, syntheticModeCtxKey, syntheticMode)
}

// IsSyntheticMode returns whether synthetic mode is enabled in the context
func IsSyntheticMode(ctx context.Context) bool {
	if mode, ok := ctx.Value(syntheticModeCtxKey).(bool); ok {
		return mode
	}
	return false
}

// GetTransactionsTable returns the appropriate transactions table name based on synthetic mode
func GetTransactionsTable(ctx context.Context) string {
	if IsSyntheticMode(ctx) {
		return "synthetic_transactions"
	}
	return "transactions"
}
