package config

import (
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
)

// WorkerConfig wraps the global config and provides worker-specific accessors
type WorkerConfig struct {
	cfg *config.Config
}

// NewWorkerConfig creates a new worker config wrapper
func NewWorkerConfig(cfg *config.Config) *WorkerConfig {
	return &WorkerConfig{cfg: cfg}
}

// Concurrency returns the number of concurrent workers
func (wc *WorkerConfig) Concurrency() int {
	return wc.cfg.Worker.Concurrency
}

// BatchSize returns the batch size for processing transactions
func (wc *WorkerConfig) BatchSize() int {
	return wc.cfg.Worker.BatchSize
}

// TransactionProcessingTimeout returns the timeout for processing a single transaction
func (wc *WorkerConfig) TransactionProcessingTimeout() time.Duration {
	return wc.cfg.Worker.Timeouts.TransactionProcessing
}

// RuleEvaluationTimeout returns the timeout for rule evaluation
func (wc *WorkerConfig) RuleEvaluationTimeout() time.Duration {
	return wc.cfg.Worker.Timeouts.RuleEvaluation
}

// RetryConfig returns the retry configuration
func (wc *WorkerConfig) RetryConfig() config.RetryConfig {
	return wc.cfg.Worker.Retry
}

// QueuePopTimeout returns the timeout for queue pop operations
func (wc *WorkerConfig) QueuePopTimeout() time.Duration {
	return wc.cfg.Worker.Queue.PopTimeout
}

// RulesReloadInterval returns the interval for reloading rules
func (wc *WorkerConfig) RulesReloadInterval() time.Duration {
	return wc.cfg.Worker.RulesReload.Interval
}
