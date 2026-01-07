package processor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewMetricsCollector_CreatesCollector(t *testing.T) {
	collector := NewMetricsCollector()

	assert.NotNil(t, collector)
	assert.NotNil(t, collector.totalProcessedCounter)
	assert.NotNil(t, collector.totalFailedCounter)
	assert.NotNil(t, collector.processingDurationHist)
}

func Test_MetricsCollector_RecordProcessing_WhenSuccess_ThenIncrementsProcessed(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordProcessing(100*time.Millisecond, true)

	metrics := collector.GetMetrics()
	assert.Equal(t, int64(1), metrics.TotalProcessed)
	assert.Equal(t, int64(0), metrics.TotalFailed)
	assert.Equal(t, 100*time.Millisecond, metrics.AverageDuration)
}

func Test_MetricsCollector_RecordProcessing_WhenFailure_ThenIncrementsFailed(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordProcessing(200*time.Millisecond, false)

	metrics := collector.GetMetrics()
	assert.Equal(t, int64(1), metrics.TotalProcessed)
	assert.Equal(t, int64(1), metrics.TotalFailed)
	assert.Equal(t, 200*time.Millisecond, metrics.AverageDuration)
}

func Test_MetricsCollector_RecordProcessing_MultipleOperations_ThenCalculatesAverage(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordProcessing(100*time.Millisecond, true)
	collector.RecordProcessing(200*time.Millisecond, true)
	collector.RecordProcessing(300*time.Millisecond, false)

	metrics := collector.GetMetrics()
	assert.Equal(t, int64(3), metrics.TotalProcessed)
	assert.Equal(t, int64(1), metrics.TotalFailed)

	expectedAvg := (100 + 200 + 300) / 3
	assert.Equal(t, time.Duration(expectedAvg)*time.Millisecond, metrics.AverageDuration)
}

func Test_MetricsCollector_GetMetrics_WhenNoData_ThenReturnsZeroValues(t *testing.T) {
	collector := NewMetricsCollector()

	metrics := collector.GetMetrics()

	assert.Equal(t, int64(0), metrics.TotalProcessed)
	assert.Equal(t, int64(0), metrics.TotalFailed)
	assert.Equal(t, time.Duration(0), metrics.AverageDuration)
	assert.Equal(t, time.Duration(0), metrics.TotalDuration)
	assert.True(t, metrics.LastProcessedTime.IsZero())
}

func Test_MetricsCollector_RecordProcessing_UpdatesLastProcessedTime(t *testing.T) {
	collector := NewMetricsCollector()

	beforeTime := time.Now()
	collector.RecordProcessing(50*time.Millisecond, true)
	afterTime := time.Now()

	metrics := collector.GetMetrics()
	assert.True(t, metrics.LastProcessedTime.After(beforeTime) || metrics.LastProcessedTime.Equal(beforeTime))
	assert.True(t, metrics.LastProcessedTime.Before(afterTime) || metrics.LastProcessedTime.Equal(afterTime))
}

func Test_MetricsCollector_RecordProcessing_Concurrent_ThenHandlesRaceConditions(t *testing.T) {
	collector := NewMetricsCollector()

	const numGoroutines = 100
	const numOpsPerGoroutine = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numOpsPerGoroutine; j++ {
				collector.RecordProcessing(10*time.Millisecond, j%2 == 0)
			}
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	metrics := collector.GetMetrics()
	expectedTotal := int64(numGoroutines * numOpsPerGoroutine)
	expectedFailed := expectedTotal / 2

	assert.Equal(t, expectedTotal, metrics.TotalProcessed)
	assert.Equal(t, expectedFailed, metrics.TotalFailed)
}

func Test_MeasureExecution_WhenFunctionSucceeds_ThenReturnsDuration(t *testing.T) {
	ctx := context.Background()

	duration, err := MeasureExecution(ctx, func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	require.NoError(t, err)
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond)
	assert.Less(t, duration, 100*time.Millisecond)
}

func Test_MeasureExecution_WhenFunctionFails_ThenReturnsError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("test error")

	duration, err := MeasureExecution(ctx, func() error {
		time.Sleep(5 * time.Millisecond)
		return expectedErr
	})

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.GreaterOrEqual(t, duration, 5*time.Millisecond)
}

func Test_MeasureBatchExecution_WhenCalled_ThenReturnsDurationAndResults(t *testing.T) {
	ctx := context.Background()

	expectedResults := []BatchResult{
		{ExternalID: "tx1", Success: true},
		{ExternalID: "tx2", Success: false},
	}

	duration, results := MeasureBatchExecution(ctx, func() []BatchResult {
		time.Sleep(10 * time.Millisecond)
		return expectedResults
	})

	assert.GreaterOrEqual(t, duration, 10*time.Millisecond)
	assert.Equal(t, expectedResults, results)
}

func Test_MeasureBatchExecution_WhenEmptyResults_ThenReturnsEmpty(t *testing.T) {
	ctx := context.Background()

	duration, results := MeasureBatchExecution(ctx, func() []BatchResult {
		time.Sleep(1 * time.Millisecond)
		return []BatchResult{}
	})

	assert.GreaterOrEqual(t, duration, 1*time.Millisecond)
	assert.Empty(t, results)
}

func Test_MetricsCollector_GetMetrics_CalculatesTotalDuration(t *testing.T) {
	collector := NewMetricsCollector()

	collector.RecordProcessing(100*time.Millisecond, true)
	collector.RecordProcessing(200*time.Millisecond, true)

	metrics := collector.GetMetrics()

	assert.Equal(t, 300*time.Millisecond, metrics.TotalDuration)
}
