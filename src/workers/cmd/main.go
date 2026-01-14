package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/database"
	"github.com/algo-shield/algo-shield/src/pkg/queue"
	"github.com/algo-shield/algo-shield/src/workers/internal/rules"
	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
)

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Info().Msg("Starting Algo Shield Worker with Asynq")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.NewPostgresPool(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Msg("Database connection established")

	// Initialize Redis (for caching)
	redisClient, err := database.NewRedisClient(cfg.GetRedisAddr())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis cache")
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing Redis connection")
		}
	}()

	log.Info().Str("redis_addr", cfg.GetRedisAddr()).Msg("Redis cache connection established")

	// Get Asynq configuration from environment
	redisQueueAddr := getEnvOrDefault("REDIS_QUEUE_HOST", "localhost") + ":" + getEnvOrDefault("REDIS_QUEUE_PORT", "6379")
	redisQueueDB, _ := strconv.Atoi(getEnvOrDefault("REDIS_QUEUE_DB", "0"))
	concurrency, _ := strconv.Atoi(getEnvOrDefault("ASYNQ_CONCURRENCY", "10"))

	queueCriticalWeight, _ := strconv.Atoi(getEnvOrDefault("ASYNQ_QUEUE_CRITICAL_WEIGHT", "6"))
	queueDefaultWeight, _ := strconv.Atoi(getEnvOrDefault("ASYNQ_QUEUE_DEFAULT_WEIGHT", "3"))
	queueLowWeight, _ := strconv.Atoi(getEnvOrDefault("ASYNQ_QUEUE_LOW_WEIGHT", "1"))

	log.Info().
		Str("redis_queue_addr", redisQueueAddr).
		Int("concurrency", concurrency).
		Msg("Asynq configuration loaded")

	// Create transaction repository
	txRepo := transactions.NewPostgresRepository(db.Pool)

	// Create rule engine
	ruleEngine := rules.NewEngine(db.Pool, redisClient.Client, cfg.Worker.Timeouts.RuleEvaluation)

	// Create transaction service (processor)
	txService := transactions.NewService(txRepo, ruleEngine)

	log.Info().Msg("Transaction service initialized")

	// Create Asynq worker configuration
	workerConfig := queue.AsynqWorkerConfig{
		RedisAddr:   redisQueueAddr,
		RedisDB:     redisQueueDB,
		Concurrency: concurrency,
		Queues: map[string]int{
			queue.QueueCritical: queueCriticalWeight,
			queue.QueueDefault:  queueDefaultWeight,
			queue.QueueLow:      queueLowWeight,
		},
		ShutdownTimeout: cfg.Worker.Asynq.ShutdownTimeout,
	}

	// Create Asynq worker
	worker, err := queue.NewAsynqWorker(workerConfig, txService)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Asynq worker")
	}

	// Register task handlers
	if err := worker.RegisterHandlers(); err != nil {
		log.Fatal().Err(err).Msg("Failed to register handlers")
	}

	log.Info().Msg("Asynq worker initialized successfully")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start worker in a goroutine
	workerErrChan := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting Asynq worker server")
		err := worker.Run()
		// Always signal completion, even if nil (graceful shutdown)
		workerErrChan <- err
	}()

	// Wait for shutdown signal or worker completion
	select {
	case sig := <-sigChan:
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case err := <-workerErrChan:
		if err != nil {
			log.Error().Err(err).Msg("Worker stopped with error")
		} else {
			log.Info().Msg("Worker stopped gracefully")
		}
	}

	// Graceful shutdown
	log.Info().Msg("Shutting down worker gracefully...")
	worker.Shutdown()
	log.Info().Msg("Worker shutdown complete")
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
