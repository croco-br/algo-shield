package schemas

import (
	"log"
	"strings"

	"github.com/expr-lang/expr"
)

// BuildExpressionEnv builds a dynamic expression environment from event JSON
// using the schema's extracted fields as the structure.
// Returns a map[string]any that can be used with expr-lang.
func BuildExpressionEnv(eventData map[string]any, schema *EventSchema) map[string]any {
	if schema == nil || eventData == nil {
		return make(map[string]any)
	}

	env := make(map[string]any)

	// For each field in the schema, extract the value from the event data
	for _, field := range schema.ExtractedFields {
		value := extractValueByPath(eventData, field.Path)
		env[field.Path] = value
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
func EvaluateExpressionWithSchema(expression string, eventData map[string]any, schema *EventSchema) bool {
	if expression == "" {
		return false
	}

	// Build expression environment from schema and event data
	env := BuildExpressionEnv(eventData, schema)

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
