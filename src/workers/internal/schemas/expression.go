package schemas

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/algo-shield/algo-shield/src/workers/internal/transactions"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// ExpressionCache caches compiled expr-lang programs for reuse
// This avoids expensive recompilation of the same expressions
type ExpressionCache struct {
	mu    sync.RWMutex
	cache map[string]*vm.Program
}

// Global expression cache instance
var expressionCache = &ExpressionCache{
	cache: make(map[string]*vm.Program),
}

// GetCompiledProgram returns a cached compiled program or compiles and caches a new one
// The cache key includes the expression and schema ID to handle different environments
func (ec *ExpressionCache) GetCompiledProgram(expression string, schemaID string, env map[string]any) (*vm.Program, error) {
	cacheKey := schemaID + ":" + expression

	// Try to get from cache (read lock)
	ec.mu.RLock()
	if program, ok := ec.cache[cacheKey]; ok {
		ec.mu.RUnlock()
		return program, nil
	}
	ec.mu.RUnlock()

	// Compile the expression (no lock during compilation)
	program, err := expr.Compile(expression, expr.Env(env), expr.AsBool())
	if err != nil {
		return nil, err
	}

	// Store in cache (write lock)
	ec.mu.Lock()
	// Double-check in case another goroutine compiled it
	if existing, ok := ec.cache[cacheKey]; ok {
		ec.mu.Unlock()
		return existing, nil
	}
	ec.cache[cacheKey] = program
	ec.mu.Unlock()

	return program, nil
}

// ClearCache clears the expression cache (useful when rules are reloaded)
func (ec *ExpressionCache) ClearCache() {
	ec.mu.Lock()
	ec.cache = make(map[string]*vm.Program)
	ec.mu.Unlock()
}

// GetExpressionCache returns the global expression cache instance
func GetExpressionCache() *ExpressionCache {
	return expressionCache
}

// BuildExpressionEnv builds a dynamic expression environment from event JSON
// using the schema's extracted fields as the structure.
// Returns a map[string]any that can be used with expr-lang.
func BuildExpressionEnv(ctx context.Context, eventData map[string]any, schema *EventSchema, historyRepo transactions.TransactionHistoryRepository) map[string]any {
	if schema == nil || eventData == nil {
		return make(map[string]any)
	}

	env := make(map[string]any)

	// For each field in the schema, extract the value and build nested structure
	// This allows expressions like "location.lat" to work via property access
	for _, field := range schema.ExtractedFields {
		value := extractValueByPath(eventData, field.Path)
		setNestedValue(env, field.Path, value)
	}

	// Add helper functions to the environment
	// The polygon parameter accepts any type because expr-lang parses array literals as []interface{}
	// We handle the type conversion internally
	env["pointInPolygon"] = func(lat, lon float64, polygon any) bool {
		converted, ok := convertToFloat64Polygon(polygon)
		if !ok {
			log.Printf("pointInPolygon: failed to convert polygon to [][]float64")
			return false
		}
		return PointInPolygon(lat, lon, converted)
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

// setNestedValue sets a value in a nested map structure using dot notation path
// e.g., setNestedValue(env, "location.lat", 37.7) creates env["location"]["lat"] = 37.7
func setNestedValue(env map[string]any, path string, value any) {
	parts := strings.Split(path, ".")

	// For single-part paths, just set directly
	if len(parts) == 1 {
		env[path] = value
		return
	}

	// Navigate/create nested maps for all but the last part
	current := env
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if existing, ok := current[part]; ok {
			// If nested map exists, use it
			if nested, ok := existing.(map[string]any); ok {
				current = nested
			} else {
				// Existing value is not a map, create new map (overwrites)
				nested := make(map[string]any)
				current[part] = nested
				current = nested
			}
		} else {
			// Create new nested map
			nested := make(map[string]any)
			current[part] = nested
			current = nested
		}
	}

	// Set the final value
	current[parts[len(parts)-1]] = value
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
// Uses expression caching for performance - expressions are compiled once and reused.
func EvaluateExpressionWithSchema(ctx context.Context, expression string, eventData map[string]any, schema *EventSchema, historyRepo transactions.TransactionHistoryRepository) bool {
	if expression == "" {
		return false
	}

	// Build expression environment from schema and event data, including helper functions
	env := BuildExpressionEnv(ctx, eventData, schema, historyRepo)

	// Get compiled program from cache or compile and cache it
	// Using schema ID as part of cache key to handle different schema environments
	schemaID := ""
	if schema != nil {
		schemaID = schema.ID.String()
	}

	program, err := expressionCache.GetCompiledProgram(expression, schemaID, env)
	if err != nil {
		log.Printf("Expression compile error: %v (expression: %s)", err, expression)
		return false
	}

	// Run the compiled expression with the current environment
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

// convertToFloat64Polygon converts an interface{} (typically []interface{} from expr-lang)
// to [][]float64 for use with PointInPolygon
func convertToFloat64Polygon(polygon any) ([][]float64, bool) {
	// If already the right type, return directly
	if typed, ok := polygon.([][]float64); ok {
		return typed, true
	}

	// Handle []interface{} from expr-lang array literals
	outerSlice, ok := polygon.([]interface{})
	if !ok {
		return nil, false
	}

	result := make([][]float64, len(outerSlice))
	for i, inner := range outerSlice {
		// Each inner element should be a slice of two floats (lat, lon)
		innerSlice, ok := inner.([]interface{})
		if !ok {
			return nil, false
		}

		if len(innerSlice) < 2 {
			return nil, false
		}

		lat, latOk := ToFloat64(innerSlice[0])
		lon, lonOk := ToFloat64(innerSlice[1])
		if !latOk || !lonOk {
			return nil, false
		}

		result[i] = []float64{lat, lon}
	}

	return result, true
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
