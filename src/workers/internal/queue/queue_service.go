package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/redis/go-redis/v9"
)

var (
	// ErrTimeout is returned when queue pop times out (expected behavior)
	ErrTimeout = errors.New("queue timeout")
	// ErrInvalidData is returned when queue data cannot be unmarshaled
	ErrInvalidData = errors.New("invalid queue data")
)

// QueueService handles transaction queue operations
type QueueService struct {
	redis *redis.Client
}

// NewQueueService creates a new queue service
func NewQueueService(redis *redis.Client) *QueueService {
	return &QueueService{redis: redis}
}

// PopTransaction pops a transaction from the Redis queue
// Returns ErrTimeout if no transaction is available (expected)
// Returns other errors for actual failures
func (q *QueueService) PopTransaction(ctx context.Context) (*models.TransactionEvent, error) {
	result, err := q.redis.BRPop(ctx, 1*time.Second, "transaction:queue").Result()

	// Check if it's a timeout (expected) vs actual error
	if err != nil {
		if err == redis.Nil {
			return nil, ErrTimeout
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, ErrInvalidData
	}

	eventJSON := result[1]
	var event models.TransactionEvent

	if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
		return nil, ErrInvalidData
	}

	return &event, nil
}
