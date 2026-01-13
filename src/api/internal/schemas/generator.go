package schemas

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EventGenerator generates synthetic events from a schema
type EventGenerator struct {
	rng *rand.Rand
}

// NewEventGenerator creates a new event generator
// Uses current time as seed
func NewEventGenerator() *EventGenerator {
	return &EventGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateEvent creates a synthetic event based on the schema's sample JSON
func (g *EventGenerator) GenerateEvent(schema *EventSchema) map[string]any {
	return g.generateFromSample(schema.SampleJSON)
}

// generateFromSample recursively generates random values based on sample structure
func (g *EventGenerator) generateFromSample(sample map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range sample {
		result[key] = g.generateValueForKey(key, value)
	}

	return result
}

// generateValueForKey generates a value considering the key name for special handling
func (g *EventGenerator) generateValueForKey(key string, value any) any {
	keyLower := strings.ToLower(key)

	// Generate UUID for external_id
	if keyLower == "external_id" || keyLower == "id" || keyLower == "event_id" {
		return uuid.New().String()
	}

	// Generate integer values for amount fields
	if keyLower == "amount" || keyLower == "value" || keyLower == "total" || keyLower == "sum" {
		var base float64
		switch v := value.(type) {
		case float64:
			base = v
		case int:
			base = float64(v)
		default:
			base = 1000
		}
		if base == 0 {
			base = 1000
		}
		// Generate integer between 0.5x and 1.5x of base value
		multiplier := 0.5 + g.rng.Float64()
		return int(base * multiplier)
	}

	// Use realistic names for origin, destination, sender, receiver, account holder fields
	if keyLower == "origin" || keyLower == "destination" || keyLower == "sender" ||
		keyLower == "receiver" || keyLower == "account_holder" || keyLower == "customer_name" ||
		keyLower == "beneficiary" || keyLower == "payer" || keyLower == "payee" {
		if _, ok := value.(string); ok {
			return g.GetRandomName()
		}
	}

	// Generate readable transaction types (debit/credit)
	if keyLower == "type" || keyLower == "transaction_type" || keyLower == "event_type" {
		types := []string{"debit", "credit"}
		return types[g.rng.Intn(len(types))]
	}

	return g.generateValue(value)
}

// generateValue generates a random value based on the sample value's type
func (g *EventGenerator) generateValue(sample any) any {
	if sample == nil {
		return nil
	}

	switch v := sample.(type) {
	case bool:
		return g.rng.Intn(2) == 1

	case float64:
		// Generate number in similar range to sample
		base := v
		if base == 0 {
			base = 1000
		}
		// Generate between 0.5x and 1.5x of base value
		multiplier := 0.5 + g.rng.Float64()
		return base * multiplier

	case int:
		base := float64(v)
		if base == 0 {
			base = 1000
		}
		multiplier := 0.5 + g.rng.Float64()
		return int(base * multiplier)

	case string:
		return g.generateString(v)

	case []any:
		return g.generateArray(v)

	case map[string]any:
		return g.generateFromSample(v)

	default:
		return sample
	}
}

// generateString generates a random string based on the sample
func (g *EventGenerator) generateString(sample string) string {
	// Check if it looks like a UUID
	if len(sample) == 36 && sample[8] == '-' && sample[13] == '-' {
		return uuid.New().String()
	}

	// Check if it looks like a date
	if _, err := time.Parse(time.RFC3339, sample); err == nil {
		// Generate a random time within the last 30 days
		days := g.rng.Intn(30)
		hours := g.rng.Intn(24)
		return time.Now().AddDate(0, 0, -days).Add(-time.Duration(hours) * time.Hour).Format(time.RFC3339)
	}

	// Check if it looks like an email
	if len(sample) > 5 && contains(sample, '@') && contains(sample, '.') {
		domains := []string{"email.com", "test.org", "example.net", "mail.io", "company.com", "finance.org"}
		// Use realistic names for email generation
		name := RealisticNames[g.rng.Intn(len(RealisticNames))]
		// Convert to email format: first.last123@domain.com
		email := strings.ToLower(strings.ReplaceAll(name, " ", "."))
		return fmt.Sprintf("%s%d@%s", email, g.rng.Intn(1000), domains[g.rng.Intn(len(domains))])
	}

	// Check if it looks like currency code (3 uppercase letters)
	if len(sample) == 3 && isUpperCase(sample) {
		currencies := []string{"USD", "EUR", "GBP", "BRL", "JPY", "CAD", "AUD"}
		return currencies[g.rng.Intn(len(currencies))]
	}

	// Check if it looks like a country code (2 uppercase letters)
	if len(sample) == 2 && isUpperCase(sample) {
		countries := []string{"US", "BR", "GB", "DE", "FR", "JP", "CA", "AU"}
		return countries[g.rng.Intn(len(countries))]
	}

	// Default: generate random alphanumeric string of similar length
	length := len(sample)
	if length == 0 {
		length = 10
	}
	return randomString(g.rng, length)
}

// GetRandomName returns a random realistic name from the names list
func (g *EventGenerator) GetRandomName() string {
	return RealisticNames[g.rng.Intn(len(RealisticNames))]
}

// generateArray generates a random array based on the sample
func (g *EventGenerator) generateArray(sample []any) []any {
	if len(sample) == 0 {
		return []any{}
	}

	// Generate 1-5 elements based on the first element's type
	count := 1 + g.rng.Intn(5)
	result := make([]any, count)

	for i := 0; i < count; i++ {
		result[i] = g.generateValue(sample[0])
	}

	return result
}

// Helper functions

func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func isUpperCase(s string) bool {
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func randomString(rng *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}
