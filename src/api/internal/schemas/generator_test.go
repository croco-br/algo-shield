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
