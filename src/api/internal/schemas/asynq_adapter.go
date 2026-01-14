package schemas

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/queue"
)

// AsynqAdapter wraps AsynqClient to implement TaskEnqueuer interface
// This adapter follows the Dependency Inversion Principle by allowing
// the schema service to work with any task enqueuer, not just Asynq
type AsynqAdapter struct {
	client *queue.AsynqClient
}

// NewAsynqAdapter creates a new adapter for AsynqClient
func NewAsynqAdapter(client *queue.AsynqClient) *AsynqAdapter {
	return &AsynqAdapter{
		client: client,
	}
}

// EnqueueTransactionWithPriority enqueues a transaction with the specified priority
// Implements TaskEnqueuer interface
func (a *AsynqAdapter) EnqueueTransactionWithPriority(ctx context.Context, event map[string]any, priority string) (any, error) {
	return a.client.EnqueueTransactionWithPriority(ctx, event, priority)
}
