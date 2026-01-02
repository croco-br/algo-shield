package processor

import (
	"context"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/workers/internal/queue"
	engine "github.com/algo-shield/algo-shield/src/workers/internal/rules"
	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Processor struct {
	transactionService  *transactions.Service
	queueService        *queue.QueueService
	ruleEngine          *engine.Engine
	metricsCollector    *MetricsCollector
	retryConfig         RetryConfig
	concurrency         int
	batchSize           int
	transactionTimeout  time.Duration
	rulesReloadInterval time.Duration
}

func NewProcessor(db *pgxpool.Pool, redis *redis.Client, concurrency, batchSize int, transactionTimeout, ruleEvaluationTimeout, queuePopTimeout, rulesReloadInterval time.Duration, retryConfig RetryConfig) *Processor {
	// Create history provider for engine
	historyProvider := transactions.NewPostgresHistoryRepository(db)

	// Create single instance of rule engine with history provider and timeout
	ruleEngine := engine.NewEngine(db, redis, historyProvider, ruleEvaluationTimeout)

	// Create transaction repository and service with dependency injection
	transactionService := transactions.NewService(transactions.NewPostgresRepository(db), ruleEngine)

	// Default batch size to 50 if not provided
	if batchSize <= 0 {
		batchSize = 50
	}

	return &Processor{
		transactionService:  transactionService,
		queueService:        queue.NewQueueService(redis, queuePopTimeout),
		ruleEngine:          ruleEngine,
		metricsCollector:    NewMetricsCollector(),
		retryConfig:         retryConfig,
		concurrency:         concurrency,
		batchSize:           batchSize,
		transactionTimeout:  transactionTimeout,
		rulesReloadInterval: rulesReloadInterval,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	log.Println("Starting transaction processor...")

	// Load rules and schemas initially
	if err := p.ruleEngine.LoadRules(ctx); err != nil {
		return err
	}

	// Create errgroup for managing all goroutines with proper error handling
	g, gCtx := errgroup.WithContext(ctx)

	// Start schema invalidation subscription (managed by errgroup)
	g.Go(func() error {
		p.ruleEngine.StartSchemaInvalidationSubscription(gCtx)
		return nil // Subscription runs until context cancellation
	})

	// Reload rules periodically for hot-reload
	g.Go(func() error {
		p.reloadRulesPeriodically(gCtx)
		return nil // Periodic reload doesn't return errors that should stop the processor
	})

	// Start worker goroutines using errgroup
	for i := 0; i < p.concurrency; i++ {
		workerID := i // Capture loop variable
		g.Go(func() error {
			p.worker(gCtx, workerID)
			return nil // Workers run until context cancellation
		})
	}

	// Wait for context cancellation or any error
	log.Println("Processor started, waiting for shutdown signal...")

	// Wait for all goroutines to finish (they'll stop when context is cancelled)
	if err := g.Wait(); err != nil {
		log.Printf("Processor stopped with error: %v", err)
		return err
	}

	log.Println("Shutdown signal received, all workers stopped gracefully")

	// Log final metrics
	metrics := p.metricsCollector.GetMetrics()
	log.Printf("Processor stopped. Metrics: processed=%d, failed=%d, avg_duration=%v",
		metrics.TotalProcessed, metrics.TotalFailed, metrics.AverageDuration)

	return nil
}

// GetMetrics returns current processing metrics
func (p *Processor) GetMetrics() Metrics {
	return p.metricsCollector.GetMetrics()
}

func (p *Processor) reloadRulesPeriodically(ctx context.Context) {
	ticker := time.NewTicker(p.rulesReloadInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.ruleEngine.LoadRules(ctx); err != nil {
				log.Printf("Failed to reload rules: %v", err)
			} else {
				log.Println("Rules reloaded successfully")
			}
		}
	}
}

func (p *Processor) worker(ctx context.Context, id int) {
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		default:
			// Process in batches if batchSize > 1, otherwise process one at a time
			if p.batchSize > 1 {
				p.processBatch(ctx)
			} else {
				p.processNextTransaction(ctx)
			}
		}
	}
}

func (p *Processor) processNextTransaction(ctx context.Context) {
	// Pop transaction from queue
	event, err := p.queueService.PopTransaction(ctx)
	if err != nil {
		// Check if it's a timeout (expected) vs actual error
		if err == queue.ErrTimeout {
			// Timeout is expected, just continue
			return
		}
		// Log actual errors
		log.Printf("Queue error: %v", err)
		return
	}

	if event == nil {
		return
	}

	// Process with metrics and retry
	duration, err := MeasureExecution(ctx, func() error {
		return Retry(ctx, p.retryConfig, func() error {
			// Add timeout to context
			processCtx, cancel := context.WithTimeout(ctx, p.transactionTimeout)
			defer cancel()

			return p.transactionService.ProcessTransaction(processCtx, *event)
		})
	})

	success := err == nil
	p.metricsCollector.RecordProcessing(duration, success)

	// Extract external_id for logging
	externalID := "unknown"
	if id, ok := (*event)["external_id"].(string); ok {
		externalID = id
	} else if id, ok := (*event)["id"].(string); ok {
		externalID = id
	}

	if err != nil {
		log.Printf("Failed to process transaction %s after retries: %v (duration: %v)", externalID, err, duration)
	} else {
		log.Printf("Processed transaction %s successfully (duration: %v)", externalID, duration)
	}
}

// processBatch processes multiple transactions in a batch using parallel processing
// Uses worker pool pattern with controlled concurrency to avoid overwhelming the system
func (p *Processor) processBatch(ctx context.Context) {
	events := make([]*models.Event, 0, p.batchSize)

	// Collect batch
	for i := 0; i < p.batchSize; i++ {
		event, err := p.queueService.PopTransaction(ctx)
		if err != nil {
			if err == queue.ErrTimeout {
				break // No more items available
			}
			log.Printf("Queue error while collecting batch: %v", err)
			continue
		}
		if event != nil {
			events = append(events, event)
		}
	}

	if len(events) == 0 {
		return
	}

	// Process batch in parallel with controlled concurrency
	duration, batchResults := MeasureBatchExecution(ctx, func() []BatchResult {
		return p.processBatchParallel(ctx, events)
	})

	// Record metrics for each transaction individually
	successCount := 0
	failureCount := 0
	for _, result := range batchResults {
		p.metricsCollector.RecordProcessing(result.Duration, result.Success)
		if result.Success {
			successCount++
		} else {
			failureCount++
			log.Printf("Failed to process transaction %s in batch: %v (duration: %v)",
				result.ExternalID, result.Error, result.Duration)
		}
	}

	// Log batch summary
	if failureCount == 0 {
		log.Printf("Processed batch of %d transactions successfully (duration: %v, avg: %v)",
			len(events), duration, duration/time.Duration(len(events)))
	} else {
		log.Printf("Processed batch of %d transactions: %d succeeded, %d failed (duration: %v)",
			len(events), successCount, failureCount, duration)
	}
}

// BatchResult represents the result of processing a single transaction in a batch
type BatchResult struct {
	ExternalID string
	Success    bool
	Error      error
	Duration   time.Duration
}

// processBatchParallel processes events in parallel with controlled concurrency
// Uses golang.org/x/sync/semaphore for robust concurrency control and errgroup for error handling
func (p *Processor) processBatchParallel(ctx context.Context, events []*models.Event) []BatchResult {
	// Use concurrency limit to prevent overwhelming the system
	// Use the processor's concurrency setting, but cap at batch size
	concurrencyLimit := int64(p.concurrency)
	if concurrencyLimit > int64(len(events)) {
		concurrencyLimit = int64(len(events))
	}
	if concurrencyLimit < 1 {
		concurrencyLimit = 1
	}

	// Create weighted semaphore for controlled concurrency
	sem := semaphore.NewWeighted(concurrencyLimit)

	// Create errgroup for managing goroutines with proper error handling
	g, gCtx := errgroup.WithContext(ctx)

	// Results channel with buffer for all events
	results := make(chan BatchResult, len(events))

	// Process each event in parallel with controlled concurrency
	for _, event := range events {
		evt := event // Capture loop variable

		g.Go(func() error {
			// Acquire semaphore (blocks if limit reached, respects context cancellation)
			if err := sem.Acquire(gCtx, 1); err != nil {
				// Context cancelled or semaphore acquisition failed
				// Extract external_id for logging
				externalID := "unknown"
				if id, ok := (*evt)["external_id"].(string); ok {
					externalID = id
				} else if id, ok := (*evt)["id"].(string); ok {
					externalID = id
				}

				results <- BatchResult{
					ExternalID: externalID,
					Success:    false,
					Error:      err,
					Duration:   0,
				}
				return err
			}
			defer sem.Release(1) // Release semaphore when done

			// Process with individual timeout and retry
			processCtx, cancel := context.WithTimeout(gCtx, p.transactionTimeout)
			defer cancel()

			duration, err := MeasureExecution(processCtx, func() error {
				return Retry(processCtx, p.retryConfig, func() error {
					return p.transactionService.ProcessTransaction(processCtx, *evt)
				})
			})

			// Extract external_id for logging
			externalID := "unknown"
			if id, ok := (*evt)["external_id"].(string); ok {
				externalID = id
			} else if id, ok := (*evt)["id"].(string); ok {
				externalID = id
			}

			results <- BatchResult{
				ExternalID: externalID,
				Success:    err == nil,
				Error:      err,
				Duration:   duration,
			}

			return nil // Don't propagate individual transaction errors to errgroup
		})
	}

	// Wait for all goroutines to complete (or context cancellation)
	// Note: We ignore the error from Wait() because we want to collect all results
	// even if some transactions failed
	_ = g.Wait()
	close(results)

	// Collect all results
	batchResults := make([]BatchResult, 0, len(events))
	for result := range results {
		batchResults = append(batchResults, result)
	}

	return batchResults
}
