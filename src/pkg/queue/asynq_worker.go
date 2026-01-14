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

// TransactionProcessor defines the interface for processing transactions
// Follows Interface Segregation Principle - worker only needs this method
type TransactionProcessor interface {
	ProcessTransaction(ctx context.Context, event models.Event) error
}

// AsynqWorkerConfig holds configuration for the Asynq worker
type AsynqWorkerConfig struct {
	RedisAddr string
	RedisDB   int

	// Concurrency: number of concurrent workers
	Concurrency int

	// Queue priorities (weight-based)
	Queues map[string]int

	// Timeouts and retries
	ShutdownTimeout time.Duration
}

// DefaultAsynqWorkerConfig returns sensible defaults for the worker
func DefaultAsynqWorkerConfig(redisAddr string) AsynqWorkerConfig {
	return AsynqWorkerConfig{
		RedisAddr:   redisAddr,
		RedisDB:     0,
		Concurrency: 10,
		Queues: map[string]int{
			QueueCritical: 6, // Higher weight = higher priority
			QueueDefault:  3,
			QueueLow:      1,
		},
		ShutdownTimeout: 30 * time.Second,
	}
}

// AsynqWorker wraps the Asynq server for consuming and processing jobs
// Follows Interface Segregation Principle (ISP) - single responsibility: consuming jobs
type AsynqWorker struct {
	server    *asynq.Server
	mux       *asynq.ServeMux
	processor TransactionProcessor
	config    AsynqWorkerConfig
}

// NewAsynqWorker creates a new Asynq worker with dependency injection
// Uses Dependency Inversion Principle - receives interface, not concrete type
func NewAsynqWorker(config AsynqWorkerConfig, processor TransactionProcessor) (*AsynqWorker, error) {
	// Validate configuration
	if config.RedisAddr == "" {
		return nil, fmt.Errorf("redis address cannot be empty")
	}
	if processor == nil {
		return nil, fmt.Errorf("transaction processor cannot be nil")
	}
	if config.Concurrency <= 0 {
		return nil, fmt.Errorf("concurrency must be positive")
	}
	if len(config.Queues) == 0 {
		return nil, fmt.Errorf("queues cannot be empty")
	}

	// Create Asynq server
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: config.RedisAddr,
			DB:   config.RedisDB,
		},
		asynq.Config{
			Concurrency: config.Concurrency,
			Queues:      config.Queues,
			// Error handler for logging failed jobs
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("task_type", task.Type()).
					Str("task_id", task.ResultWriter().TaskID()).
					Msg("Task execution failed")
			}),
			// Graceful shutdown timeout
			ShutdownTimeout: config.ShutdownTimeout,
			// Health check function
			HealthCheckFunc: func(err error) {
				if err != nil {
					log.Error().Err(err).Msg("Health check failed")
				}
			},
		},
	)

	mux := asynq.NewServeMux()

	log.Info().
		Str("redis_addr", config.RedisAddr).
		Int("concurrency", config.Concurrency).
		Interface("queues", config.Queues).
		Msg("Asynq worker initialized")

	return &AsynqWorker{
		server:    server,
		mux:       mux,
		processor: processor,
		config:    config,
	}, nil
}

// RegisterHandlers registers all task handlers
// Must be called before Run()
func (w *AsynqWorker) RegisterHandlers() error {
	if w.mux == nil {
		return fmt.Errorf("mux is nil")
	}

	// Register transaction processing handler
	w.mux.HandleFunc(TaskTypeProcessTransaction, w.handleProcessTransaction)

	log.Info().
		Str("task_type", TaskTypeProcessTransaction).
		Msg("Registered task handler")

	return nil
}

// handleProcessTransaction is the handler for processing transaction events
// Deserializes the event and delegates to the transaction processor
func (w *AsynqWorker) handleProcessTransaction(ctx context.Context, task *asynq.Task) error {
	startTime := time.Now()

	// Extract task metadata for logging
	taskID := task.ResultWriter().TaskID()

	log.Info().
		Str("task_id", taskID).
		Msg("Processing transaction task")

	// Deserialize event from payload
	var event models.Event
	if err := json.Unmarshal(task.Payload(), &event); err != nil {
		// Fatal error: payload is corrupted, no point in retrying
		log.Error().
			Err(err).
			Str("task_id", taskID).
			Msg("Failed to unmarshal event payload - marking as failed")
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Extract external_id for logging
	externalID := extractExternalID(event)

	// Process transaction using injected processor
	if err := w.processor.ProcessTransaction(ctx, event); err != nil {
		// Log error and return for Asynq to handle retry
		log.Error().
			Err(err).
			Str("task_id", taskID).
			Str("external_id", externalID).
			Msg("Failed to process transaction")
		return fmt.Errorf("failed to process transaction: %w", err)
	}

	// Calculate processing time
	processingTime := time.Since(startTime)

	log.Info().
		Str("task_id", taskID).
		Str("external_id", externalID).
		Int64("processing_time_ms", processingTime.Milliseconds()).
		Msg("Transaction processed successfully")

	return nil
}

// Run starts the Asynq worker server (blocking call)
// Should be called in a goroutine if non-blocking behavior is needed
func (w *AsynqWorker) Run() error {
	if w.server == nil {
		return fmt.Errorf("server is nil")
	}
	if w.mux == nil {
		return fmt.Errorf("mux is nil")
	}

	log.Info().Msg("Starting Asynq worker server")

	// This is a blocking call that runs until Shutdown() is called
	if err := w.server.Run(w.mux); err != nil {
		log.Error().Err(err).Msg("Asynq worker server stopped with error")
		return fmt.Errorf("worker server error: %w", err)
	}

	log.Info().Msg("Asynq worker server stopped gracefully")
	return nil
}

// Shutdown gracefully shuts down the Asynq worker
// Waits for in-flight tasks to complete (up to ShutdownTimeout)
func (w *AsynqWorker) Shutdown() {
	if w.server == nil {
		log.Warn().Msg("Server is nil, nothing to shutdown")
		return
	}

	log.Info().
		Float64("timeout_seconds", w.config.ShutdownTimeout.Seconds()).
		Msg("Shutting down Asynq worker gracefully")

	// Shutdown waits for in-flight tasks to complete
	w.server.Shutdown()

	log.Info().Msg("Asynq worker shutdown complete")
}

// Health checks if the Asynq worker is healthy
// Returns error if the server is not running or Redis is unreachable
func (w *AsynqWorker) Health(ctx context.Context) error {
	if w.server == nil {
		return fmt.Errorf("worker server is nil")
	}

	// Try to ping Redis through an inspector
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: w.config.RedisAddr,
		DB:   w.config.RedisDB,
	})
	defer func() {
		if closeErr := inspector.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("Failed to close inspector")
		}
	}()

	// Try to get queue info (lightweight health check)
	queues, err := inspector.Queues()
	if err != nil {
		log.Error().Err(err).Msg("Worker health check failed")
		return fmt.Errorf("worker health check failed: %w", err)
	}

	log.Debug().
		Strs("queues", queues).
		Msg("Worker health check passed")

	return nil
}
