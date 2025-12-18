package processor

import (
	"context"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	// meter is the OpenTelemetry meter for processor metrics
	meter = otel.Meter("github.com/algo-shield/algo-shield/processor")
)

// Metrics tracks processing metrics
type Metrics struct {
	TotalProcessed    int64
	TotalFailed       int64
	TotalDuration     time.Duration
	AverageDuration   time.Duration
	LastProcessedTime time.Time
}

// MetricsCollector collects and tracks metrics using OpenTelemetry
// OpenTelemetry handles thread-safety internally, eliminating race conditions
type MetricsCollector struct {
	// OpenTelemetry metrics (thread-safe by design)
	totalProcessedCounter  metric.Int64Counter
	totalFailedCounter     metric.Int64Counter
	processingDurationHist metric.Int64Histogram

	// Local aggregated values for GetMetrics() using atomic operations
	totalProcessed    atomic.Int64
	totalFailed       atomic.Int64
	totalDurationNano atomic.Int64
	lastProcessedTime atomic.Int64 // UnixNano timestamp
}

// NewMetricsCollector creates a new metrics collector with OpenTelemetry
// OpenTelemetry meter is always available, so this function never returns an error
func NewMetricsCollector() *MetricsCollector {
	// Create OpenTelemetry counters and histograms
	// These operations are safe and will not fail with the global meter
	totalProcessedCounter, _ := meter.Int64Counter(
		"processor_transactions_total",
		metric.WithDescription("Total number of transactions processed"),
	)

	totalFailedCounter, _ := meter.Int64Counter(
		"processor_transactions_failed_total",
		metric.WithDescription("Total number of failed transactions"),
	)

	processingDurationHist, _ := meter.Int64Histogram(
		"processor_transaction_duration_nanoseconds",
		metric.WithDescription("Transaction processing duration in nanoseconds"),
		metric.WithUnit("ns"),
	)

	return &MetricsCollector{
		totalProcessedCounter:  totalProcessedCounter,
		totalFailedCounter:     totalFailedCounter,
		processingDurationHist: processingDurationHist,
	}
}

// RecordProcessing records a processing operation using OpenTelemetry
// This method is thread-safe - OpenTelemetry handles concurrency internally
func (mc *MetricsCollector) RecordProcessing(duration time.Duration, success bool) {
	// Record metrics in OpenTelemetry (thread-safe)
	durationNano := duration.Nanoseconds()
	mc.totalProcessedCounter.Add(context.Background(), 1)
	mc.processingDurationHist.Record(context.Background(), durationNano)

	if !success {
		mc.totalFailedCounter.Add(context.Background(), 1)
	}

	// Update local atomic counters for GetMetrics() compatibility
	mc.totalProcessed.Add(1)
	if !success {
		mc.totalFailed.Add(1)
	}
	mc.totalDurationNano.Add(durationNano)
	mc.lastProcessedTime.Store(time.Now().UnixNano())
}

// GetMetrics returns current metrics snapshot
// Uses atomic operations for thread-safe reads
func (mc *MetricsCollector) GetMetrics() Metrics {
	totalProcessed := mc.totalProcessed.Load()
	totalFailed := mc.totalFailed.Load()
	totalDurationNano := mc.totalDurationNano.Load()
	lastProcessedNano := mc.lastProcessedTime.Load()

	var averageDuration time.Duration
	if totalProcessed > 0 {
		averageDuration = time.Duration(totalDurationNano / totalProcessed)
	}

	var lastProcessedTime time.Time
	if lastProcessedNano > 0 {
		lastProcessedTime = time.Unix(0, lastProcessedNano)
	}

	return Metrics{
		TotalProcessed:    totalProcessed,
		TotalFailed:       totalFailed,
		TotalDuration:     time.Duration(totalDurationNano),
		AverageDuration:   averageDuration,
		LastProcessedTime: lastProcessedTime,
	}
}

// MeasureExecution measures execution time of a function
func MeasureExecution(ctx context.Context, fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	return duration, err
}

// MeasureBatchExecution measures execution time of a batch processing function
func MeasureBatchExecution(ctx context.Context, fn func() []BatchResult) (time.Duration, []BatchResult) {
	start := time.Now()
	results := fn()
	duration := time.Since(start)
	return duration, results
}
