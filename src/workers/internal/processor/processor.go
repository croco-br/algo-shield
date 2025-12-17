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
)

type Processor struct {
	transactionService *transactions.Service
	queueService       *queue.QueueService
	ruleEngine         *engine.Engine
	metricsCollector   *MetricsCollector
	retryConfig        RetryConfig
	concurrency        int
	batchSize          int
}

func NewProcessor(db *pgxpool.Pool, redis *redis.Client, concurrency, batchSize int) *Processor {
	// Create history provider for engine
	historyProvider := transactions.NewPostgresHistoryRepository(db)

	// Create single instance of rule engine with history provider
	ruleEngine := engine.NewEngine(db, redis, historyProvider)

	// Create transaction service with injected engine
	transactionService := transactions.NewService(db, redis, ruleEngine)

	// Default batch size to 50 if not provided
	if batchSize <= 0 {
		batchSize = 50
	}

	return &Processor{
		transactionService: transactionService,
		queueService:       queue.NewQueueService(redis),
		ruleEngine:         ruleEngine,
		metricsCollector:   NewMetricsCollector(),
		retryConfig:        DefaultRetryConfig(),
		concurrency:        concurrency,
		batchSize:          batchSize,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	log.Println("Starting transaction processor...")

	// Load rules initially
	if err := p.ruleEngine.LoadRules(ctx); err != nil {
		return err
	}

	// Reload rules periodically for hot-reload
	go p.reloadRulesPeriodically(ctx)

	// Start worker goroutines
	for i := 0; i < p.concurrency; i++ {
		go p.worker(ctx, i)
	}

	<-ctx.Done()

	// Log final metrics
	metrics := p.metricsCollector.GetMetrics()
	log.Printf("Processor stopping. Metrics: processed=%d, failed=%d, avg_duration=%v",
		metrics.TotalProcessed, metrics.TotalFailed, metrics.AverageDuration)

	log.Println("Stopping transaction processor...")
	return nil
}

// GetMetrics returns current processing metrics
func (p *Processor) GetMetrics() Metrics {
	return p.metricsCollector.GetMetrics()
}

func (p *Processor) reloadRulesPeriodically(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
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
			processCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
			defer cancel()

			return p.transactionService.ProcessTransaction(processCtx, *event)
		})
	})

	success := err == nil
	p.metricsCollector.RecordProcessing(duration, success)

	if err != nil {
		log.Printf("Failed to process transaction %s after retries: %v (duration: %v)", event.ExternalID, err, duration)
	} else {
		log.Printf("Processed transaction %s successfully (duration: %v)", event.ExternalID, duration)
	}
}

// processBatch processes multiple transactions in a batch
func (p *Processor) processBatch(ctx context.Context) {
	events := make([]*models.TransactionEvent, 0, p.batchSize)

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

	// Process batch with metrics
	duration, err := MeasureExecution(ctx, func() error {
		batchCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond*time.Duration(len(events)))
		defer cancel()

		for _, event := range events {
			if err := Retry(batchCtx, p.retryConfig, func() error {
				return p.transactionService.ProcessTransaction(batchCtx, *event)
			}); err != nil {
				return err
			}
		}
		return nil
	})

	success := err == nil
	p.metricsCollector.RecordProcessing(duration, success)

	if err != nil {
		log.Printf("Failed to process batch of %d transactions: %v (duration: %v)", len(events), err, duration)
	} else {
		log.Printf("Processed batch of %d transactions successfully (duration: %v, avg: %v)",
			len(events), duration, duration/time.Duration(len(events)))
	}
}
