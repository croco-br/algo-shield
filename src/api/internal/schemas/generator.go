package schemas

import (
	"crypto/rand"
	"encoding/binary"
	"strings"
	"time"
)

// EventGenerator generates synthetic events from a schema
type EventGenerator struct{}

// NewEventGenerator creates a new event generator
func NewEventGenerator() *EventGenerator {
	return &EventGenerator{}
}

// GenerateEvent creates a synthetic event based on the schema's extracted fields
func (g *EventGenerator) GenerateEvent(schema *EventSchema) map[string]any {
	result := make(map[string]any)

	// Use extracted_fields to build the event structure
	for _, field := range schema.ExtractedFields {
		value := g.generateValueFromField(&field)
		g.setNestedValue(result, field.Path, value)
	}

	// If no extracted fields, fallback to sample JSON
	if len(result) == 0 {
		return g.generateFromSample(schema.SampleJSON)
	}

	return result
}

// generateValueFromField generates a value based solely on the extracted field type
func (g *EventGenerator) generateValueFromField(field *ExtractedField) any {
	switch field.Type {
	case FieldTypeString:
		// Check if field path suggests it's a date/time field
		if g.isDateField(field.Path) {
			return g.generateRandomDate()
		}
		return g.generateRandomString()
	case FieldTypeDateTime:
		return g.generateRandomDate()
	case FieldTypeNumber:
		return g.generateRandomNumber()
	case FieldTypeBoolean:
		return cryptoRandomInt(2) == 1
	case FieldTypeArray:
		return []any{}
	case FieldTypeObject:
		return map[string]any{}
	case FieldTypeNull:
		return nil
	default:
		return g.generateRandomString()
	}
}

// isDateField checks if a field path suggests it's a date/time field
func (g *EventGenerator) isDateField(path string) bool {
	pathLower := strings.ToLower(path)
	dateKeywords := []string{"date", "time", "timestamp", "created", "updated", "at", "when", "since"}
	for _, keyword := range dateKeywords {
		if strings.Contains(pathLower, keyword) {
			return true
		}
	}
	return false
}

// generateRandomDate generates a random RFC3339 date string
func (g *EventGenerator) generateRandomDate() string {
	// Generate a random time within the last 90 days
	daysAgo := cryptoRandomInt(90)
	hoursAgo := cryptoRandomInt(24)
	minutesAgo := cryptoRandomInt(60)

	return time.Now().
		AddDate(0, 0, -daysAgo).
		Add(-time.Duration(hoursAgo) * time.Hour).
		Add(-time.Duration(minutesAgo) * time.Minute).
		Format(time.RFC3339)
}

// generateRandomString generates a random string
func (g *EventGenerator) generateRandomString() string {
	return randomString(10 + cryptoRandomInt(20))
}

// generateRandomNumber generates a random number
func (g *EventGenerator) generateRandomNumber() float64 {
	return float64(cryptoRandomInt(10000))
}

// setNestedValue sets a value in a nested map using dot notation path
func (g *EventGenerator) setNestedValue(m map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	current := m

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, set the value
			current[part] = value
		} else {
			// Not the last part, ensure nested map exists
			if _, ok := current[part]; !ok {
				current[part] = make(map[string]any)
			}
			if nested, ok := current[part].(map[string]any); ok {
				current = nested
			} else {
				// Path conflict, create new nested map
				nested := make(map[string]any)
				current[part] = nested
				current = nested
			}
		}
	}
}

// generateFromSample recursively generates random values based on sample structure
// Used only as fallback when no extracted_fields are available
func (g *EventGenerator) generateFromSample(sample map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range sample {
		result[key] = g.generateValueByType(value)
	}

	return result
}

// generateValueByType generates a random value based solely on the value's type
func (g *EventGenerator) generateValueByType(sample any) any {
	if sample == nil {
		return nil
	}

	switch v := sample.(type) {
	case bool:
		return cryptoRandomInt(2) == 1

	case float64:
		return float64(cryptoRandomInt(10000))

	case int:
		return cryptoRandomInt(10000)

	case string:
		return randomString(10 + cryptoRandomInt(20))

	case []any:
		return g.generateArray(v)

	case map[string]any:
		return g.generateFromSample(v)

	default:
		return randomString(10)
	}
}

// generateArray generates a random array
func (g *EventGenerator) generateArray(sample []any) []any {
	if len(sample) == 0 {
		return []any{}
	}

	// Generate 1-3 elements based on the first element's type
	count := 1 + cryptoRandomInt(3)
	result := make([]any, count)

	for i := 0; i < count; i++ {
		result[i] = g.generateValueByType(sample[0])
	}

	return result
}

// Helper functions

// cryptoRandomInt generates a cryptographically secure random integer in the range [0, max)
func cryptoRandomInt(max int) int {
	if max <= 0 {
		return 0
	}

	// Calculate the number of bytes needed
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// Fallback to 0 if crypto/rand fails (should never happen)
		return 0
	}

	// Convert bytes to uint64
	n := binary.BigEndian.Uint64(b[:])

	// Use modulo to get value in range [0, max)
	// This introduces slight bias for non-power-of-2 values, but acceptable for synthetic data
	return int(n % uint64(max))
}

// randomString generates a cryptographically secure random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[cryptoRandomInt(len(charset))]
	}
	return string(b)
}
