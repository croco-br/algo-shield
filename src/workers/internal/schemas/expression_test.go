package schemas

import (
	"context"
	"sync"
	"testing"

	"github.com/expr-lang/expr/vm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExpressionCache Tests

func Test_ExpressionCache_GetCompiledProgram_WhenFirstCall_ThenCompilesAndCaches(t *testing.T) {
	cache := &ExpressionCache{cache: make(map[string]*vm.Program)}
	expression := "amount > 100"
	schemaID := uuid.New().String()
	env := map[string]any{"amount": 150.0}

	program, err := cache.GetCompiledProgram(expression, schemaID, env)

	require.NoError(t, err)
	assert.NotNil(t, program)
	assert.Len(t, cache.cache, 1)
}

func Test_ExpressionCache_GetCompiledProgram_WhenCacheHit_ThenReturnsCachedProgram(t *testing.T) {
	cache := &ExpressionCache{cache: make(map[string]*vm.Program)}
	expression := "amount > 100"
	schemaID := uuid.New().String()
	env := map[string]any{"amount": 150.0}

	program1, err1 := cache.GetCompiledProgram(expression, schemaID, env)
	program2, err2 := cache.GetCompiledProgram(expression, schemaID, env)

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Same(t, program1, program2)
	assert.Len(t, cache.cache, 1)
}

func Test_ExpressionCache_GetCompiledProgram_WhenInvalidExpression_ThenReturnsError(t *testing.T) {
	cache := &ExpressionCache{cache: make(map[string]*vm.Program)}
	expression := "invalid syntax &&& ||"
	schemaID := uuid.New().String()
	env := map[string]any{}

	program, err := cache.GetCompiledProgram(expression, schemaID, env)

	assert.Error(t, err)
	assert.Nil(t, program)
}

func Test_ExpressionCache_GetCompiledProgram_WhenConcurrentAccess_ThenHandlesSafely(t *testing.T) {
	cache := &ExpressionCache{cache: make(map[string]*vm.Program)}
	expression := "amount > 100"
	schemaID := uuid.New().String()
	env := map[string]any{"amount": 150.0}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			program, err := cache.GetCompiledProgram(expression, schemaID, env)
			assert.NoError(t, err)
			assert.NotNil(t, program)
		}()
	}

	wg.Wait()
	assert.Len(t, cache.cache, 1)
}

func Test_ExpressionCache_ClearCache_WhenCalled_ThenRemovesAllEntries(t *testing.T) {
	cache := &ExpressionCache{cache: make(map[string]*vm.Program)}
	_, _ = cache.GetCompiledProgram("amount > 100", "schema1", map[string]any{"amount": 150.0})
	_, _ = cache.GetCompiledProgram("value > 200", "schema2", map[string]any{"value": 250.0})

	cache.ClearCache()

	assert.Empty(t, cache.cache)
}

// BuildExpressionEnv Tests

func Test_BuildExpressionEnv_WhenNilSchema_ThenReturnsEmptyMap(t *testing.T) {
	eventData := map[string]any{"amount": 100.0}

	env := BuildExpressionEnv(context.Background(), eventData, nil, nil)

	assert.Empty(t, env)
}

func Test_BuildExpressionEnv_WhenNilEventData_ThenReturnsEmptyMap(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
		},
	}

	env := BuildExpressionEnv(context.Background(), nil, schema, nil)

	assert.Empty(t, env)
}

func Test_BuildExpressionEnv_WhenSimpleFields_ThenExtractsCorrectly(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
			{Path: "currency", Type: FieldTypeString},
		},
	}
	eventData := map[string]any{
		"amount":   100.0,
		"currency": "USD",
	}

	env := BuildExpressionEnv(context.Background(), eventData, schema, nil)

	assert.Equal(t, 100.0, env["amount"])
	assert.Equal(t, "USD", env["currency"])
}

func Test_BuildExpressionEnv_WhenNestedFields_ThenExtractsCorrectly(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "user.id", Type: FieldTypeString},
			{Path: "user.country", Type: FieldTypeString},
		},
	}
	eventData := map[string]any{
		"user": map[string]any{
			"id":      "user123",
			"country": "BR",
		},
	}

	env := BuildExpressionEnv(context.Background(), eventData, schema, nil)

	userMap, ok := env["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "user123", userMap["id"])
	assert.Equal(t, "BR", userMap["country"])
}

func Test_BuildExpressionEnv_WhenHistoryRepoProvided_ThenIncludesHelperFunctions(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "origin", Type: FieldTypeString},
		},
	}
	eventData := map[string]any{"origin": "192.168.1.1"}
	mockRepo := &mockHistoryRepository{}

	env := BuildExpressionEnv(context.Background(), eventData, schema, mockRepo)

	assert.Contains(t, env, "velocityCount")
	assert.Contains(t, env, "velocitySum")
	assert.Contains(t, env, "pointInPolygon")
}

// EvaluateExpressionWithSchema Tests

func Test_EvaluateExpressionWithSchema_WhenEmptyExpression_ThenReturnsFalse(t *testing.T) {
	schema := &EventSchema{ID: uuid.New()}
	eventData := map[string]any{"amount": 100.0}

	result := EvaluateExpressionWithSchema(context.Background(), "", eventData, schema, nil)

	assert.False(t, result)
}

func Test_EvaluateExpressionWithSchema_WhenValidExpressionTrue_ThenReturnsTrue(t *testing.T) {
	schema := &EventSchema{
		ID: uuid.New(),
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
		},
	}
	eventData := map[string]any{"amount": 150.0}

	result := EvaluateExpressionWithSchema(context.Background(), "amount > 100", eventData, schema, nil)

	assert.True(t, result)
}

func Test_EvaluateExpressionWithSchema_WhenValidExpressionFalse_ThenReturnsFalse(t *testing.T) {
	schema := &EventSchema{
		ID: uuid.New(),
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
		},
	}
	eventData := map[string]any{"amount": 50.0}

	result := EvaluateExpressionWithSchema(context.Background(), "amount > 100", eventData, schema, nil)

	assert.False(t, result)
}

func Test_EvaluateExpressionWithSchema_WhenCompilationError_ThenReturnsFalse(t *testing.T) {
	schema := &EventSchema{ID: uuid.New()}
	eventData := map[string]any{"amount": 100.0}

	result := EvaluateExpressionWithSchema(context.Background(), "invalid &&& syntax", eventData, schema, nil)

	assert.False(t, result)
}

func Test_EvaluateExpressionWithSchema_WhenExpressionUsesCache_ThenReusesCachedProgram(t *testing.T) {
	cache := GetExpressionCache()
	cache.ClearCache()
	schema := &EventSchema{
		ID: uuid.New(),
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
		},
	}
	eventData1 := map[string]any{"amount": 150.0}
	eventData2 := map[string]any{"amount": 50.0}

	result1 := EvaluateExpressionWithSchema(context.Background(), "amount > 100", eventData1, schema, nil)
	result2 := EvaluateExpressionWithSchema(context.Background(), "amount > 100", eventData2, schema, nil)

	assert.True(t, result1)
	assert.False(t, result2)
}

// PointInPolygon Tests

func Test_PointInPolygon_WhenPointInside_ThenReturnsTrue(t *testing.T) {
	polygon := [][]float64{
		{0, 0},
		{0, 10},
		{10, 10},
		{10, 0},
	}

	result := PointInPolygon(5, 5, polygon)

	assert.True(t, result)
}

func Test_PointInPolygon_WhenPointOutside_ThenReturnsFalse(t *testing.T) {
	polygon := [][]float64{
		{0, 0},
		{0, 10},
		{10, 10},
		{10, 0},
	}

	result := PointInPolygon(15, 15, polygon)

	assert.False(t, result)
}

func Test_PointInPolygon_WhenPointOnEdge_ThenReturnsConsistentResult(t *testing.T) {
	polygon := [][]float64{
		{0, 0},
		{0, 10},
		{10, 10},
		{10, 0},
	}

	result := PointInPolygon(0, 5, polygon)

	assert.True(t, result)
}

func Test_PointInPolygon_WhenPolygonTooSmall_ThenReturnsFalse(t *testing.T) {
	polygon := [][]float64{
		{0, 0},
		{10, 10},
	}

	result := PointInPolygon(5, 5, polygon)

	assert.False(t, result)
}

// ToFloat64 Tests

func Test_ToFloat64_WhenFloat64_ThenReturnsValue(t *testing.T) {
	result, ok := ToFloat64(42.5)

	assert.True(t, ok)
	assert.Equal(t, 42.5, result)
}

func Test_ToFloat64_WhenFloat32_ThenConvertsToFloat64(t *testing.T) {
	result, ok := ToFloat64(float32(42.5))

	assert.True(t, ok)
	assert.InDelta(t, 42.5, result, 0.001)
}

func Test_ToFloat64_WhenInt_ThenConvertsToFloat64(t *testing.T) {
	result, ok := ToFloat64(42)

	assert.True(t, ok)
	assert.Equal(t, 42.0, result)
}

func Test_ToFloat64_WhenInt64_ThenConvertsToFloat64(t *testing.T) {
	result, ok := ToFloat64(int64(42))

	assert.True(t, ok)
	assert.Equal(t, 42.0, result)
}

func Test_ToFloat64_WhenInt32_ThenConvertsToFloat64(t *testing.T) {
	result, ok := ToFloat64(int32(42))

	assert.True(t, ok)
	assert.Equal(t, 42.0, result)
}

func Test_ToFloat64_WhenString_ThenReturnsFalse(t *testing.T) {
	result, ok := ToFloat64("42")

	assert.False(t, ok)
	assert.Equal(t, 0.0, result)
}

// convertToFloat64Polygon Tests

func Test_ConvertToFloat64Polygon_WhenAlreadyCorrectType_ThenReturnsDirectly(t *testing.T) {
	polygon := [][]float64{{0, 0}, {10, 10}, {20, 0}}

	result, ok := convertToFloat64Polygon(polygon)

	assert.True(t, ok)
	assert.Equal(t, polygon, result)
}

func Test_ConvertToFloat64Polygon_WhenInterfaceSlice_ThenConverts(t *testing.T) {
	polygon := []interface{}{
		[]interface{}{0.0, 0.0},
		[]interface{}{10.0, 10.0},
		[]interface{}{20.0, 0.0},
	}

	result, ok := convertToFloat64Polygon(polygon)

	assert.True(t, ok)
	assert.Equal(t, [][]float64{{0, 0}, {10, 10}, {20, 0}}, result)
}

func Test_ConvertToFloat64Polygon_WhenInvalidType_ThenReturnsFalse(t *testing.T) {
	polygon := "not a polygon"

	result, ok := convertToFloat64Polygon(polygon)

	assert.False(t, ok)
	assert.Nil(t, result)
}

func Test_ConvertToFloat64Polygon_WhenInvalidInnerType_ThenReturnsFalse(t *testing.T) {
	polygon := []interface{}{
		"invalid",
		[]interface{}{10.0, 10.0},
	}

	result, ok := convertToFloat64Polygon(polygon)

	assert.False(t, ok)
	assert.Nil(t, result)
}

func Test_ConvertToFloat64Polygon_WhenInnerSliceTooSmall_ThenReturnsFalse(t *testing.T) {
	polygon := []interface{}{
		[]interface{}{0.0},
	}

	result, ok := convertToFloat64Polygon(polygon)

	assert.False(t, ok)
	assert.Nil(t, result)
}

// setNestedValue Tests

func Test_SetNestedValue_WhenSinglePath_ThenSetsDirectly(t *testing.T) {
	env := make(map[string]any)

	setNestedValue(env, "amount", 100.0)

	assert.Equal(t, 100.0, env["amount"])
}

func Test_SetNestedValue_WhenNestedPath_ThenCreatesNesting(t *testing.T) {
	env := make(map[string]any)

	setNestedValue(env, "user.id", "user123")

	userMap, ok := env["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "user123", userMap["id"])
}

func Test_SetNestedValue_WhenDeeplyNestedPath_ThenCreatesAllLevels(t *testing.T) {
	env := make(map[string]any)

	setNestedValue(env, "location.address.city", "São Paulo")

	locationMap, ok := env["location"].(map[string]any)
	require.True(t, ok)
	addressMap, ok := locationMap["address"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "São Paulo", addressMap["city"])
}

func Test_SetNestedValue_WhenPathAlreadyExists_ThenOverwrites(t *testing.T) {
	env := map[string]any{
		"user": map[string]any{
			"id": "old_id",
		},
	}

	setNestedValue(env, "user.id", "new_id")

	userMap, ok := env["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "new_id", userMap["id"])
}

// extractValueByPath Tests

func Test_ExtractValueByPath_WhenNilData_ThenReturnsNil(t *testing.T) {
	result := extractValueByPath(nil, "amount")

	assert.Nil(t, result)
}

func Test_ExtractValueByPath_WhenSinglePath_ThenReturnsValue(t *testing.T) {
	data := map[string]any{"amount": 100.0}

	result := extractValueByPath(data, "amount")

	assert.Equal(t, 100.0, result)
}

func Test_ExtractValueByPath_WhenNestedPath_ThenReturnsNestedValue(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"id": "user123",
		},
	}

	result := extractValueByPath(data, "user.id")

	assert.Equal(t, "user123", result)
}

func Test_ExtractValueByPath_WhenPathNotExists_ThenReturnsNil(t *testing.T) {
	data := map[string]any{"amount": 100.0}

	result := extractValueByPath(data, "nonexistent")

	assert.Nil(t, result)
}

func Test_ExtractValueByPath_WhenNestedPathNotExists_ThenReturnsNil(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"id": "user123",
		},
	}

	result := extractValueByPath(data, "user.nonexistent")

	assert.Nil(t, result)
}

// detectSumField Tests

func Test_DetectSumField_WhenNilSchema_ThenReturnsFallback(t *testing.T) {
	result := detectSumField(nil)

	assert.Equal(t, "amount", result)
}

func Test_DetectSumField_WhenAmountExists_ThenReturnsAmount(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
			{Path: "value", Type: FieldTypeNumber},
		},
	}

	result := detectSumField(schema)

	assert.Equal(t, "amount", result)
}

func Test_DetectSumField_WhenOnlyValueExists_ThenReturnsValue(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "value", Type: FieldTypeNumber},
			{Path: "total", Type: FieldTypeNumber},
		},
	}

	result := detectSumField(schema)

	assert.Equal(t, "value", result)
}

func Test_DetectSumField_WhenOnlyTotalExists_ThenReturnsTotal(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "total", Type: FieldTypeNumber},
		},
	}

	result := detectSumField(schema)

	assert.Equal(t, "total", result)
}

func Test_DetectSumField_WhenNoPreferredFields_ThenReturnsFirstNumericField(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "name", Type: FieldTypeString},
			{Path: "price", Type: FieldTypeNumber},
		},
	}

	result := detectSumField(schema)

	assert.Equal(t, "price", result)
}

func Test_DetectSumField_WhenNoNumericFields_ThenReturnsFallback(t *testing.T) {
	schema := &EventSchema{
		ExtractedFields: []ExtractedField{
			{Path: "name", Type: FieldTypeString},
		},
	}

	result := detectSumField(schema)

	assert.Equal(t, "amount", result)
}

// extractFieldPathAndValue Tests

func Test_ExtractFieldPathAndValue_WhenValidPath_ThenReturnsPathAndValue(t *testing.T) {
	eventData := map[string]any{"origin": "192.168.1.1"}
	schema := &EventSchema{}

	path, value := extractFieldPathAndValue("origin", eventData, schema)

	assert.Equal(t, "origin", path)
	assert.Equal(t, "192.168.1.1", value)
}

func Test_ExtractFieldPathAndValue_WhenNotString_ThenReturnsEmpty(t *testing.T) {
	eventData := map[string]any{"origin": "192.168.1.1"}
	schema := &EventSchema{}

	path, value := extractFieldPathAndValue(123, eventData, schema)

	assert.Empty(t, path)
	assert.Empty(t, value)
}

func Test_ExtractFieldPathAndValue_WhenPathNotExists_ThenReturnsEmpty(t *testing.T) {
	eventData := map[string]any{"origin": "192.168.1.1"}
	schema := &EventSchema{}

	path, value := extractFieldPathAndValue("nonexistent", eventData, schema)

	assert.Empty(t, path)
	assert.Empty(t, value)
}

func Test_ExtractFieldPathAndValue_WhenNumericValue_ThenConvertsToString(t *testing.T) {
	eventData := map[string]any{"user_id": 12345}
	schema := &EventSchema{}

	path, value := extractFieldPathAndValue("user_id", eventData, schema)

	assert.Equal(t, "user_id", path)
	assert.Equal(t, "12345", value)
}

// Mock implementations

type mockHistoryRepository struct {
	countResult int
	countErr    error
	sumResult   float64
	sumErr      error
}

func (m *mockHistoryRepository) CountByFieldInTimeWindow(ctx context.Context, groupFieldPath string, fieldValue string, timeWindowSeconds int) (int, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return m.countResult, nil
}

func (m *mockHistoryRepository) SumFieldByFieldInTimeWindow(ctx context.Context, groupFieldPath string, groupFieldValue string, sumFieldPath string, timeWindowSeconds int) (float64, error) {
	if m.sumErr != nil {
		return 0.0, m.sumErr
	}
	return m.sumResult, nil
}
