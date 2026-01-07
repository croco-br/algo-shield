package processor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DefaultRetryConfig_ReturnsValidConfig(t *testing.T) {
	config := DefaultRetryConfig()

	assert.Equal(t, 3, config.MaxAttempts)
	assert.Equal(t, 100*time.Millisecond, config.InitialDelay)
	assert.Equal(t, 5*time.Second, config.MaxDelay)
	assert.Equal(t, 2.0, config.Multiplier)
}

func Test_Retry_WhenFunctionSucceeds_ThenReturnsNoError(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	callCount := 0
	fn := func() error {
		callCount++
		return nil
	}

	err := Retry(ctx, config, fn)

	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func Test_Retry_WhenFunctionFailsOnce_ThenRetries(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	callCount := 0
	fn := func() error {
		callCount++
		if callCount < 2 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := Retry(ctx, config, fn)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func Test_Retry_WhenFunctionAlwaysFails_ThenReturnsLastError(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	expectedErr := errors.New("permanent error")
	callCount := 0
	fn := func() error {
		callCount++
		return expectedErr
	}

	err := Retry(ctx, config, fn)

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 3, callCount)
}

func Test_Retry_WhenContextCancelled_ThenReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}

	fn := func() error {
		return errors.New("should not be called")
	}

	err := Retry(ctx, config, fn)

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func Test_Retry_WhenContextCancelledDuringDelay_ThenReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}

	callCount := 0
	fn := func() error {
		callCount++
		if callCount == 1 {
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
		}
		return errors.New("error")
	}

	err := Retry(ctx, config, fn)

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.LessOrEqual(t, callCount, 2)
}

func Test_Retry_WithRetryableError_ThenUsesCustomDelay(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}

	callCount := 0
	customDelay := 5 * time.Millisecond
	fn := func() error {
		callCount++
		if callCount < 2 {
			return &RetryableError{
				Err:        errors.New("retryable error"),
				RetryAfter: customDelay,
			}
		}
		return nil
	}

	start := time.Now()
	err := Retry(ctx, config, fn)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
	assert.GreaterOrEqual(t, elapsed, customDelay)
}

func Test_RetryableError_Error_ReturnsWrappedError(t *testing.T) {
	innerErr := errors.New("inner error")
	retryErr := &RetryableError{
		Err:        innerErr,
		RetryAfter: 1 * time.Second,
	}

	assert.Equal(t, "inner error", retryErr.Error())
}

func Test_Retry_ExponentialBackoff_ThenDelaysIncrease(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  4,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}

	callCount := 0
	var delays []time.Duration
	lastCallTime := time.Now()

	fn := func() error {
		callCount++
		now := time.Now()
		if callCount > 1 {
			delays = append(delays, now.Sub(lastCallTime))
		}
		lastCallTime = now
		return errors.New("error")
	}

	err := Retry(ctx, config, fn)

	require.Error(t, err)
	assert.Equal(t, 4, callCount)
	assert.Len(t, delays, 3)

	for i := 1; i < len(delays); i++ {
		assert.GreaterOrEqual(t, delays[i], delays[i-1])
	}
}

func Test_Retry_MaxDelay_ThenCapsDelay(t *testing.T) {
	ctx := context.Background()
	config := RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     150 * time.Millisecond,
		Multiplier:   3.0,
	}

	callCount := 0
	var delays []time.Duration
	lastCallTime := time.Now()

	fn := func() error {
		callCount++
		now := time.Now()
		if callCount > 1 {
			delays = append(delays, now.Sub(lastCallTime))
		}
		lastCallTime = now
		return errors.New("error")
	}

	err := Retry(ctx, config, fn)

	require.Error(t, err)

	for _, delay := range delays {
		assert.LessOrEqual(t, delay, 200*time.Millisecond)
	}
}
