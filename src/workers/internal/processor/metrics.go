package processor

import (
	"context"
	"sync"
	"time"
)

// Metrics tracks processing metrics
type Metrics struct {
	TotalProcessed    int64
	TotalFailed       int64
	TotalDuration     time.Duration
	AverageDuration   time.Duration
	LastProcessedTime time.Time
}

// MetricsCollector collects and tracks metrics
type MetricsCollector struct {
	mu      sync.RWMutex
	metrics Metrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: Metrics{},
	}
}

// RecordProcessing records a processing operation
func (mc *MetricsCollector) RecordProcessing(duration time.Duration, success bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.metrics.TotalProcessed++
	if !success {
		mc.metrics.TotalFailed++
	}
	mc.metrics.TotalDuration += duration
	mc.metrics.LastProcessedTime = time.Now()

	// Calculate average
	if mc.metrics.TotalProcessed > 0 {
		mc.metrics.AverageDuration = mc.metrics.TotalDuration / time.Duration(mc.metrics.TotalProcessed)
	}
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() Metrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.metrics
}

// MeasureExecution measures execution time of a function
func MeasureExecution(ctx context.Context, fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	return duration, err
}
