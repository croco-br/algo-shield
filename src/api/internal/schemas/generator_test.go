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
