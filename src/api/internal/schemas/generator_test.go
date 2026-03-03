package schemas

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EventGenerator_GenerateEvent_WhenDateTimeField_ThenGeneratesDateTime(t *testing.T) {
	generator := NewEventGenerator()
	schema := &EventSchema{
		ID:   uuid.New(),
		Name: "test-schema",
		ExtractedFields: []ExtractedField{
			{Path: "created_at", Type: FieldTypeDateTime, Nullable: false},
			{Path: "updated_at", Type: FieldTypeDateTime, Nullable: false},
		},
	}

	event := generator.GenerateEvent(schema)

	require.NotNil(t, event)
	createdAt, ok := event["created_at"].(string)
	require.True(t, ok, "created_at should be a string")
	_, err := time.Parse(time.RFC3339, createdAt)
	assert.NoError(t, err, "created_at should be a valid RFC3339 datetime")

	updatedAt, ok := event["updated_at"].(string)
	require.True(t, ok, "updated_at should be a string")
	_, err = time.Parse(time.RFC3339, updatedAt)
	assert.NoError(t, err, "updated_at should be a valid RFC3339 datetime")
}

func Test_EventGenerator_GenerateEvent_WhenMixedFields_ThenGeneratesCorrectTypes(t *testing.T) {
	generator := NewEventGenerator()
	schema := &EventSchema{
		ID:   uuid.New(),
		Name: "test-schema",
		ExtractedFields: []ExtractedField{
			{Path: "name", Type: FieldTypeString, Nullable: false},
			{Path: "age", Type: FieldTypeNumber, Nullable: false},
			{Path: "active", Type: FieldTypeBoolean, Nullable: false},
			{Path: "created_at", Type: FieldTypeDateTime, Nullable: false},
		},
	}

	event := generator.GenerateEvent(schema)

	require.NotNil(t, event)
	assert.IsType(t, "", event["name"])
	assert.IsType(t, float64(0), event["age"])
	assert.IsType(t, false, event["active"])

	createdAt, ok := event["created_at"].(string)
	require.True(t, ok)
	_, err := time.Parse(time.RFC3339, createdAt)
	assert.NoError(t, err)
}

func Test_EventGenerator_GenerateEvent_WhenNestedDateTimeField_ThenGeneratesNestedDateTime(t *testing.T) {
	generator := NewEventGenerator()
	schema := &EventSchema{
		ID:   uuid.New(),
		Name: "test-schema",
		ExtractedFields: []ExtractedField{
			{Path: "user.created_at", Type: FieldTypeDateTime, Nullable: false},
			{Path: "user.name", Type: FieldTypeString, Nullable: false},
		},
	}

	event := generator.GenerateEvent(schema)

	require.NotNil(t, event)
	user, ok := event["user"].(map[string]any)
	require.True(t, ok, "user should be a map")

	createdAt, ok := user["created_at"].(string)
	require.True(t, ok, "user.created_at should be a string")
	_, err := time.Parse(time.RFC3339, createdAt)
	assert.NoError(t, err, "user.created_at should be a valid RFC3339 datetime")
}

func Test_EventGenerator_GenerateValueFromField_WhenDateTimeType_ThenReturnsDateTimeString(t *testing.T) {
	generator := NewEventGenerator()
	field := &ExtractedField{
		Path: "timestamp",
		Type: FieldTypeDateTime,
	}

	value := generator.generateValueFromField(field)

	valueStr, ok := value.(string)
	require.True(t, ok, "value should be a string")
	_, err := time.Parse(time.RFC3339, valueStr)
	assert.NoError(t, err, "value should be a valid RFC3339 datetime")
}

func Test_EventGenerator_GenerateRandomString_ThenGeneratesDifferentValues(t *testing.T) {
	generator := NewEventGenerator()

	// Generate multiple random strings to verify they're different
	values := make(map[string]bool)
	for i := 0; i < 50; i++ {
		str := generator.generateRandomString()
		values[str] = true
	}

	// With crypto/rand, all 50 strings should be unique
	assert.Equal(t, 50, len(values), "all generated strings should be unique")
}

func Test_EventGenerator_GenerateRandomNumber_ThenGeneratesDifferentValues(t *testing.T) {
	generator := NewEventGenerator()

	// Generate multiple random numbers to verify they're different
	values := make(map[float64]bool)
	for i := 0; i < 50; i++ {
		num := generator.generateRandomNumber()
		values[num] = true
	}

	// With crypto/rand, we expect high diversity (at least 45 unique out of 50)
	assert.GreaterOrEqual(t, len(values), 45, "should generate diverse random numbers")
}

func Test_cryptoRandomInt_ThenGeneratesDiverseValues(t *testing.T) {
	// Test that cryptoRandomInt produces diverse values across multiple calls
	values := make(map[int]bool)
	max := 100

	for i := 0; i < 100; i++ {
		val := cryptoRandomInt(max)
		assert.GreaterOrEqual(t, val, 0, "value should be >= 0")
		assert.Less(t, val, max, "value should be < max")
		values[val] = true
	}

	// With crypto/rand and 100 values in range [0,100), expect at least 50 unique
	// (This is a reasonable expectation given birthday paradox statistics)
	assert.GreaterOrEqual(t, len(values), 50, "should generate diverse random integers")
}

func Test_cryptoRandomInt_WhenZeroOrNegative_ThenReturnsZero(t *testing.T) {
	// Test edge cases with zero and negative max values
	assert.Equal(t, 0, cryptoRandomInt(0), "should return 0 for max=0")
	assert.Equal(t, 0, cryptoRandomInt(-1), "should return 0 for negative max")
	assert.Equal(t, 0, cryptoRandomInt(-100), "should return 0 for negative max")
}

func Test_cryptoRandomInt_WhenLargeMax_ThenReturnsValidRange(t *testing.T) {
	// Test with large max values to ensure no overflow
	largeMax := 1000000

	for i := 0; i < 100; i++ {
		val := cryptoRandomInt(largeMax)
		assert.GreaterOrEqual(t, val, 0, "value should be >= 0")
		assert.Less(t, val, largeMax, "value should be < max")
	}
}

func Test_cryptoRandomInt_WhenSmallMax_ThenReturnsValidRange(t *testing.T) {
	// Test with small max values (1, 2, 3)
	for max := 1; max <= 3; max++ {
		for i := 0; i < 50; i++ {
			val := cryptoRandomInt(max)
			assert.GreaterOrEqual(t, val, 0, "value should be >= 0")
			assert.Less(t, val, max, "value should be < max")
		}
	}
}

func Test_EventGenerator_GenerateFromSample_WhenFlatMap_ThenGeneratesRandomValues(t *testing.T) {
	generator := NewEventGenerator()
	sample := map[string]any{
		"name":   "John",
		"age":    30.0,
		"active": true,
	}

	result := generator.generateFromSample(sample)

	require.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.IsType(t, "", result["name"])
	assert.IsType(t, float64(0), result["age"])
	assert.IsType(t, false, result["active"])
}

func Test_EventGenerator_GenerateFromSample_WhenNestedMap_ThenGeneratesNestedValues(t *testing.T) {
	generator := NewEventGenerator()
	sample := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  25.0,
		},
	}

	result := generator.generateFromSample(sample)

	require.NotNil(t, result)
	user, ok := result["user"].(map[string]any)
	require.True(t, ok, "user should be a nested map")
	assert.IsType(t, "", user["name"])
	assert.IsType(t, float64(0), user["age"])
}

func Test_EventGenerator_GenerateFromSample_WhenEmptyMap_ThenReturnsEmptyMap(t *testing.T) {
	generator := NewEventGenerator()
	sample := map[string]any{}

	result := generator.generateFromSample(sample)

	require.NotNil(t, result)
	assert.Empty(t, result)
}

func Test_EventGenerator_GenerateValueByType_WhenNil_ThenReturnsNil(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType(nil)

	assert.Nil(t, result)
}

func Test_EventGenerator_GenerateValueByType_WhenBool_ThenReturnsBool(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType(true)

	assert.IsType(t, false, result)
}

func Test_EventGenerator_GenerateValueByType_WhenFloat64_ThenReturnsFloat64(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType(99.5)

	assert.IsType(t, float64(0), result)
}

func Test_EventGenerator_GenerateValueByType_WhenInt_ThenReturnsInt(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType(42)

	assert.IsType(t, 0, result)
}

func Test_EventGenerator_GenerateValueByType_WhenString_ThenReturnsString(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType("hello")

	assert.IsType(t, "", result)
}

func Test_EventGenerator_GenerateValueByType_WhenSlice_ThenReturnsSlice(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType([]any{"a", "b"})

	arr, ok := result.([]any)
	require.True(t, ok, "result should be a slice")
	assert.NotEmpty(t, arr)
	// Elements should be strings since the sample element is a string
	for _, v := range arr {
		assert.IsType(t, "", v)
	}
}

func Test_EventGenerator_GenerateValueByType_WhenMap_ThenReturnsMap(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateValueByType(map[string]any{"key": "value"})

	m, ok := result.(map[string]any)
	require.True(t, ok, "result should be a map")
	assert.Contains(t, m, "key")
}

func Test_EventGenerator_GenerateValueByType_WhenUnknownType_ThenReturnsString(t *testing.T) {
	generator := NewEventGenerator()

	// Pass something that doesn't match any known type (e.g., a struct)
	type custom struct{}
	result := generator.generateValueByType(custom{})

	assert.IsType(t, "", result)
}

func Test_EventGenerator_GenerateArray_WhenEmptyArray_ThenReturnsEmptyArray(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateArray([]any{})

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func Test_EventGenerator_GenerateArray_WhenStringElements_ThenReturnsStringArray(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateArray([]any{"hello"})

	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
	assert.LessOrEqual(t, len(result), 3)
	for _, v := range result {
		assert.IsType(t, "", v)
	}
}

func Test_EventGenerator_GenerateArray_WhenNumericElements_ThenReturnsNumericArray(t *testing.T) {
	generator := NewEventGenerator()

	result := generator.generateArray([]any{42.0})

	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result), 1)
	for _, v := range result {
		assert.IsType(t, float64(0), v)
	}
}

func Test_EventGenerator_GenerateEvent_WhenNoExtractedFields_ThenFallsBackToSample(t *testing.T) {
	generator := NewEventGenerator()
	schema := &EventSchema{
		ID:              uuid.New(),
		Name:            "test-schema",
		ExtractedFields: []ExtractedField{}, // Empty fields
		SampleJSON: map[string]any{
			"amount":   100.50,
			"currency": "USD",
		},
	}

	event := generator.GenerateEvent(schema)

	require.NotNil(t, event)
	assert.Contains(t, event, "amount")
	assert.Contains(t, event, "currency")
	assert.IsType(t, float64(0), event["amount"])
	assert.IsType(t, "", event["currency"])
}
