package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithSyntheticMode_WhenEnabled_ThenStoresInContext(t *testing.T) {
	ctx := context.Background()

	newCtx := WithSyntheticMode(ctx, true)

	assert.True(t, IsSyntheticMode(newCtx))
}

func Test_WithSyntheticMode_WhenDisabled_ThenStoresInContext(t *testing.T) {
	ctx := context.Background()

	newCtx := WithSyntheticMode(ctx, false)

	assert.False(t, IsSyntheticMode(newCtx))
}

func Test_IsSyntheticMode_WhenNotSet_ThenReturnsFalse(t *testing.T) {
	ctx := context.Background()

	result := IsSyntheticMode(ctx)

	assert.False(t, result)
}

func Test_IsSyntheticMode_WhenSetToTrue_ThenReturnsTrue(t *testing.T) {
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, true)

	result := IsSyntheticMode(ctx)

	assert.True(t, result)
}

func Test_IsSyntheticMode_WhenSetToFalse_ThenReturnsFalse(t *testing.T) {
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, false)

	result := IsSyntheticMode(ctx)

	assert.False(t, result)
}

func Test_GetTransactionsTable_WhenSyntheticMode_ThenReturnsSyntheticTable(t *testing.T) {
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, true)

	tableName := GetTransactionsTable(ctx)

	assert.Equal(t, "synthetic_transactions", tableName)
}

func Test_GetTransactionsTable_WhenNormalMode_ThenReturnsNormalTable(t *testing.T) {
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, false)

	tableName := GetTransactionsTable(ctx)

	assert.Equal(t, "transactions", tableName)
}

func Test_GetTransactionsTable_WhenNotSet_ThenReturnsNormalTable(t *testing.T) {
	ctx := context.Background()

	tableName := GetTransactionsTable(ctx)

	assert.Equal(t, "transactions", tableName)
}

func Test_WithSyntheticMode_WhenCalledMultipleTimes_ThenOverwritesValue(t *testing.T) {
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, true)
	assert.True(t, IsSyntheticMode(ctx))

	ctx = WithSyntheticMode(ctx, false)

	assert.False(t, IsSyntheticMode(ctx))
}

func Test_WithSyntheticMode_WhenChainedContexts_ThenPreservesValue(t *testing.T) {
	type otherKey struct{}
	ctx := context.Background()
	ctx = WithSyntheticMode(ctx, true)
	childCtx := context.WithValue(ctx, otherKey{}, "other_value")

	result := IsSyntheticMode(childCtx)

	assert.True(t, result)
}
