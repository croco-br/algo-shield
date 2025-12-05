package processor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/pkg/models"
	"github.com/algo-shield/algo-shield/workers/internal/rules"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Processor struct {
	db          *pgxpool.Pool
	redis       *redis.Client
	engine      *rules.Engine
	concurrency int
	batchSize   int
}

func NewProcessor(db *pgxpool.Pool, redis *redis.Client, concurrency, batchSize int) *Processor {
	return &Processor{
		db:          db,
		redis:       redis,
		engine:      rules.NewEngine(db, redis),
		concurrency: concurrency,
		batchSize:   batchSize,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	log.Println("Starting transaction processor...")

	// Load rules initially
	if err := p.engine.LoadRules(ctx); err != nil {
		return err
	}

	// Reload rules periodically for hot-reload
	go p.reloadRulesPeriodically(ctx)

	// Start worker goroutines
	for i := 0; i < p.concurrency; i++ {
		go p.worker(ctx, i)
	}

	<-ctx.Done()
	log.Println("Stopping transaction processor...")
	return nil
}

func (p *Processor) reloadRulesPeriodically(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.engine.LoadRules(ctx); err != nil {
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
			p.processNextTransaction(ctx)
		}
	}
}

func (p *Processor) processNextTransaction(ctx context.Context) {
	// Block and wait for transaction from queue
	result, err := p.redis.BRPop(ctx, 1*time.Second, "transaction:queue").Result()
	if err != nil {
		return // Timeout or error, continue
	}

	if len(result) < 2 {
		return
	}

	eventJSON := result[1]
	var event models.TransactionEvent

	if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
		log.Printf("Failed to unmarshal transaction event: %v", err)
		return
	}

	// Process transaction
	if err := p.processTransaction(ctx, event); err != nil {
		log.Printf("Failed to process transaction: %v", err)
	}
}

func (p *Processor) processTransaction(ctx context.Context, event models.TransactionEvent) error {
	startTime := time.Now()

	// Evaluate transaction against rules
	result, err := p.engine.Evaluate(ctx, event)
	if err != nil {
		return err
	}

	// Store transaction in database
	transactionID := uuid.New()
	now := time.Now()
	processingTime := time.Since(startTime).Milliseconds()

	matchedRulesJSON, _ := json.Marshal(result.MatchedRules)
	metadataJSON, _ := json.Marshal(event.Metadata)

	query := `
		INSERT INTO transactions (
			id, external_id, amount, currency, from_account, to_account, 
			type, status, risk_score, risk_level, processing_time, 
			matched_rules, metadata, created_at, processed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = p.db.Exec(ctx, query,
		transactionID, event.ExternalID, event.Amount, event.Currency,
		event.FromAccount, event.ToAccount, event.Type, result.Status,
		result.RiskScore, result.RiskLevel, processingTime,
		matchedRulesJSON, metadataJSON, now, now,
	)

	if err != nil {
		return err
	}

	log.Printf(
		"Processed transaction %s: status=%s, risk_score=%.2f, risk_level=%s, time=%dms",
		event.ExternalID, result.Status, result.RiskScore, result.RiskLevel, processingTime,
	)

	return nil
}

