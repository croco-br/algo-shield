package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/database"
	"github.com/algo-shield/algo-shield/src/workers/internal/processor"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresPool(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redis, err := database.NewRedisClient(cfg.GetRedisAddr())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Create processor
	proc := processor.NewProcessor(
		db.Pool,
		redis.Client,
		cfg.Worker.Concurrency,
		cfg.Worker.BatchSize,
	)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal...")
		cancel()
	}()

	// Start processor
	if err := proc.Start(ctx); err != nil {
		log.Fatalf("Processor failed: %v", err)
	}
}
