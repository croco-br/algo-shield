package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	"github.com/algo-shield/algo-shield/src/pkg/models"
)

const (
	// TaskTypeProcessTransaction is the task type for transaction processing
	TaskTypeProcessTransaction = "transaction:process"

	// Queue names
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"

	// Default values
	defaultMaxRetry  = 3
	defaultTimeout   = 5 * time.Minute
	defaultRetention = 24 * time.Hour
)

// AsynqClient wraps the Asynq client for enqueueing transaction jobs
// Follows Interface Segregation Principle (ISP) - single responsibility: enqueueing jobs
type AsynqClient struct {
	client *asynq.Client
}

// NewAsynqClient creates a new Asynq client with the provided Redis address
// Uses dependency injection pattern for better testability
func NewAsynqClient(redisAddr string) (*AsynqClient, error) {
	if redisAddr == "" {
		return nil, fmt.Errorf("redis address cannot be empty")
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   0,
	})

	log.Info().
		Str("redis_addr", redisAddr).
		Msg("Asynq client initialized")

	return &AsynqClient{
		client: client,
	}, nil
}

// EnqueueTransactionOpts configures transaction job enqueueing
type EnqueueTransactionOpts struct {
	Queue     string        // Queue name (critical, default, low)
	MaxRetry  int           // Maximum retry attempts
	Timeout   time.Duration // Processing timeout
	Retention time.Duration // How long to keep completed jobs
	ProcessIn time.Duration // Delay before processing (scheduled jobs)
}

// DefaultEnqueueOpts returns sensible defaults following project conventions
func DefaultEnqueueOpts() EnqueueTransactionOpts {
	return EnqueueTransactionOpts{
		Queue:     QueueDefault,
		MaxRetry:  defaultMaxRetry,
		Timeout:   defaultTimeout,
		Retention: defaultRetention,
		ProcessIn: 0, // Immediate processing
	}
}

// EnqueueTransaction enqueues a transaction event for processing
// Returns task info with job ID for tracking
// Follows fail-fast principle with comprehensive error handling
func (c *AsynqClient) EnqueueTransaction(ctx context.Context, event models.Event, opts *EnqueueTransactionOpts) (*asynq.TaskInfo, error) {
	// Validation: event cannot be empty
	if len(event) == 0 {
		return nil, fmt.Errorf("event cannot be empty")
	}

	// Use default options if not provided
	if opts == nil {
		defaultOpts := DefaultEnqueueOpts()
		opts = &defaultOpts
	}

	// Validate options
	if opts.MaxRetry < 0 {
		return nil, fmt.Errorf("max retry cannot be negative")
	}
	if opts.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}

	// Serialize event to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		log.Error().
			Err(err).
			Interface("event", event).
			Msg("Failed to marshal event")
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Asynq task
	task := asynq.NewTask(TaskTypeProcessTransaction, payload)

	// Enqueue with options
	info, err := c.client.EnqueueContext(ctx, task,
		asynq.Queue(opts.Queue),
		asynq.MaxRetry(opts.MaxRetry),
		asynq.Timeout(opts.Timeout),
		asynq.Retention(opts.Retention),
		asynq.ProcessIn(opts.ProcessIn),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("queue", opts.Queue).
			Int("max_retry", opts.MaxRetry).
			Msg("Failed to enqueue task")
		return nil, fmt.Errorf("failed to enqueue task: %w", err)
	}

	// Extract external_id for logging (if present)
	externalID := extractExternalID(event)

	log.Info().
		Str("task_id", info.ID).
		Str("queue", info.Queue).
		Str("external_id", externalID).
		Int("max_retry", info.MaxRetry).
		Msg("Transaction enqueued successfully")

	return info, nil
}

// EnqueueTransactionWithPriority is a helper for priority-based enqueueing
// Automatically selects queue based on priority level
func (c *AsynqClient) EnqueueTransactionWithPriority(ctx context.Context, event models.Event, priority string) (*asynq.TaskInfo, error) {
	opts := DefaultEnqueueOpts()

	// Map priority to queue and timeout
	switch priority {
	case "critical":
		opts.Queue = QueueCritical
		opts.Timeout = 30 * time.Second // Lower timeout for critical
	case "low":
		opts.Queue = QueueLow
		opts.Timeout = 10 * time.Minute // Higher timeout for low priority
	default:
		opts.Queue = QueueDefault
	}

	return c.EnqueueTransaction(ctx, event, &opts)
}

// Close closes the Asynq client connection
// Should be called during graceful shutdown
func (c *AsynqClient) Close() error {
	if c.client == nil {
		return nil
	}

	err := c.client.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close Asynq client")
		return fmt.Errorf("failed to close asynq client: %w", err)
	}

	log.Info().Msg("Asynq client closed")
	return nil
}

// Health checks if the Asynq client can reach Redis
// Returns error if connection fails
func (c *AsynqClient) Health(ctx context.Context, redisAddr string) error {
	if c.client == nil {
		return fmt.Errorf("asynq client is nil")
	}

	// Try to get queue info (lightweight health check)
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: redisAddr,
	})
	defer func() {
		if closeErr := inspector.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("Failed to close inspector")
		}
	}()

	_, err := inspector.Queues()
	if err != nil {
		log.Error().Err(err).Msg("Asynq health check failed")
		return fmt.Errorf("asynq health check failed: %w", err)
	}

	return nil
}

// extractExternalID safely extracts external_id from event for logging
// Returns "unknown" if not found, avoiding nil panics
func extractExternalID(event models.Event) string {
	if event == nil {
		return "unknown"
	}

	// Try multiple common field names (flexible field extraction)
	for _, key := range []string{"external_id", "id", "event_id", "transaction_id"} {
		if id, ok := event[key].(string); ok && id != "" {
			return id
		}
	}

	return "unknown"
}
