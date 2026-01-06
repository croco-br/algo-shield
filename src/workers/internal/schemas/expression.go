package schemas

import (
	"context"
	"log"
	"strings"

	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/expr-lang/expr"
)

// BuildExpressionEnv builds a dynamic expression environment from event JSON
// using the schema's extracted fields as the structure.
// Returns a map[string]any that can be used with expr-lang.
func BuildExpressionEnv(ctx context.Context, eventData map[string]any, schema *EventSchema, historyRepo transactions.TransactionHistoryRepository) map[string]any {
	if schema == nil || eventData == nil {
		return make(map[string]any)
	}

	env := make(map[string]any)

	// For each field in the schema, extract the value from the event data
	for _, field := range schema.ExtractedFields {
		value := extractValueByPath(eventData, field.Path)
		env[field.Path] = value
	}

	// Add helper functions to the environment
	env["pointInPolygon"] = func(lat, lon float64, polygon [][]float64) bool {
		return PointInPolygon(lat, lon, polygon)
	}

	// Add velocity helper functions if history repository is available
	if historyRepo != nil {
		env["velocityCount"] = func(account string, timeWindowSeconds int) int {
			count, err := historyRepo.CountByAccountInTimeWindow(ctx, account, timeWindowSeconds)
			if err != nil {
				log.Printf("Velocity count error: %v", err)
				return 0
			}
			return count
		}

		env["velocitySum"] = func(account string, timeWindowSeconds int) float64 {
			sum, err := historyRepo.SumAmountByAccountInTimeWindow(ctx, account, timeWindowSeconds)
			if err != nil {
				log.Printf("Velocity sum error: %v", err)
				return 0.0
			}
			return sum
		}
	}

	return env
}

// extractValueByPath extracts a value from nested JSON using dot notation
// e.g., "user.country" extracts data["user"]["country"]
func extractValueByPath(data map[string]any, path string) any {
	if data == nil {
		return nil
	}

	parts := strings.Split(path, ".")
	var current any = data

	for _, part := range parts {
		if current == nil {
			return nil
		}

		switch v := current.(type) {
		case map[string]any:
			current = v[part]
		default:
			return nil
		}
	}

	return current
}

// EvaluateExpressionWithSchema compiles and evaluates an expression against event data
// using a schema-defined environment.
// Returns true if the expression evaluates to true, false otherwise.
func EvaluateExpressionWithSchema(ctx context.Context, expression string, eventData map[string]any, schema *EventSchema, historyRepo transactions.TransactionHistoryRepository) bool {
	if expression == "" {
		return false
	}

	// Build expression environment from schema and event data, including helper functions
	env := BuildExpressionEnv(ctx, eventData, schema, historyRepo)

	// Compile the expression with type safety
	// expr.AsBool() ensures the result must be a boolean
	program, err := expr.Compile(expression, expr.Env(env), expr.AsBool())
	if err != nil {
		log.Printf("Expression compile error: %v (expression: %s)", err, expression)
		return false
	}

	// Run the compiled expression
	result, err := expr.Run(program, env)
	if err != nil {
		log.Printf("Expression runtime error: %v (expression: %s)", err, expression)
		return false
	}

	// The result should be a boolean due to expr.AsBool() option
	if boolResult, ok := result.(bool); ok {
		return boolResult
	}

	log.Printf("Expression did not return boolean: %T (expression: %s)", result, expression)
	return false
}

// ToFloat64 converts various numeric types to float64
func ToFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}

// PointInPolygon checks if a point is inside a polygon using the ray casting algorithm.
// This is a standard algorithm that counts how many times a ray from the point
// crosses the polygon boundary. If the count is odd, the point is inside.
func PointInPolygon(lat, lon float64, polygon [][]float64) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}

	inside := false
	j := n - 1

	for i := 0; i < n; i++ {
		// Check if the ray from (lat, lon) going right crosses the edge from polygon[i] to polygon[j]
		// polygon[i] and polygon[j] are [lat, lon] pairs
		latI, lonI := polygon[i][0], polygon[i][1]
		latJ, lonJ := polygon[j][0], polygon[j][1]

		if ((lonI > lon) != (lonJ > lon)) &&
			(lat < (latJ-latI)*(lon-lonI)/(lonJ-lonI)+latI) {
			inside = !inside
		}
		j = i
	}

	return inside
}
