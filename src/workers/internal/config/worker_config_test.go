package config

import (
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_NewWorkerConfig_WhenCalled_ThenReturnsConfig(t *testing.T) {
	cfg := &config.Config{}

	workerCfg := NewWorkerConfig(cfg)

	assert.NotNil(t, workerCfg)
}

func Test_Concurrency_WhenConfigured_ThenReturnsValue(t *testing.T) {
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Concurrency: 5,
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.Concurrency()

	assert.Equal(t, 5, result)
}

func Test_BatchSize_WhenConfigured_ThenReturnsValue(t *testing.T) {
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			BatchSize: 50,
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.BatchSize()

	assert.Equal(t, 50, result)
}

func Test_TransactionProcessingTimeout_WhenConfigured_ThenReturnsValue(t *testing.T) {
	expectedTimeout := 300 * time.Millisecond
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Timeouts: config.WorkerTimeouts{
				TransactionProcessing: expectedTimeout,
			},
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.TransactionProcessingTimeout()

	assert.Equal(t, expectedTimeout, result)
}

func Test_RuleEvaluationTimeout_WhenConfigured_ThenReturnsValue(t *testing.T) {
	expectedTimeout := 200 * time.Millisecond
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Timeouts: config.WorkerTimeouts{
				RuleEvaluation: expectedTimeout,
			},
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.RuleEvaluationTimeout()

	assert.Equal(t, expectedTimeout, result)
}

func Test_RetryConfig_WhenConfigured_ThenReturnsValue(t *testing.T) {
	expectedRetry := config.RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Second,
		MaxDelay:     10 * time.Second,
	}
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Retry: expectedRetry,
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.RetryConfig()

	assert.Equal(t, expectedRetry, result)
	assert.Equal(t, 3, result.MaxAttempts)
	assert.Equal(t, time.Second, result.InitialDelay)
	assert.Equal(t, 10*time.Second, result.MaxDelay)
}

func Test_QueuePopTimeout_WhenConfigured_ThenReturnsValue(t *testing.T) {
	expectedTimeout := 5 * time.Second
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Queue: config.QueueConfig{
				PopTimeout: expectedTimeout,
			},
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.QueuePopTimeout()

	assert.Equal(t, expectedTimeout, result)
}

func Test_RulesReloadInterval_WhenConfigured_ThenReturnsValue(t *testing.T) {
	expectedInterval := 10 * time.Second
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			RulesReload: config.RulesReloadConfig{
				Interval: expectedInterval,
			},
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	result := workerCfg.RulesReloadInterval()

	assert.Equal(t, expectedInterval, result)
}

func Test_AllAccessors_WhenCalledMultipleTimes_ThenReturnConsistentValues(t *testing.T) {
	cfg := &config.Config{
		Worker: config.WorkerConfig{
			Concurrency: 10,
			BatchSize:   100,
			Timeouts: config.WorkerTimeouts{
				TransactionProcessing: 300 * time.Millisecond,
				RuleEvaluation:        200 * time.Millisecond,
			},
			Retry: config.RetryConfig{
				MaxAttempts:  3,
				InitialDelay: time.Second,
			},
			Queue: config.QueueConfig{
				PopTimeout: 5 * time.Second,
			},
			RulesReload: config.RulesReloadConfig{
				Interval: 10 * time.Second,
			},
		},
	}
	workerCfg := NewWorkerConfig(cfg)

	assert.Equal(t, 10, workerCfg.Concurrency())
	assert.Equal(t, 10, workerCfg.Concurrency())
	assert.Equal(t, 100, workerCfg.BatchSize())
	assert.Equal(t, 100, workerCfg.BatchSize())
	assert.Equal(t, 300*time.Millisecond, workerCfg.TransactionProcessingTimeout())
	assert.Equal(t, 200*time.Millisecond, workerCfg.RuleEvaluationTimeout())
}
