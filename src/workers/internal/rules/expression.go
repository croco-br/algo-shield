package rules

import (
	"log"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/expr-lang/expr"
)

// ExprEnv is the environment passed to expr for custom rule evaluation.
// Contains all transaction fields that can be referenced in expressions.
// Field names use snake_case to match the expression syntax users will write.
type ExprEnv struct {
	// Transaction fields (directly from TransactionEvent)
	Amount      float64        `expr:"amount"`
	Currency    string         `expr:"currency"`
	Type        string         `expr:"type"`
	FromAccount string         `expr:"from_account"`
	ToAccount   string         `expr:"to_account"`
	ExternalID  string         `expr:"external_id"`
	Metadata    map[string]any `expr:"metadata"`

	// Convenience fields from metadata (flattened for easier access)
	Region    string  `expr:"region"`
	Country   string  `expr:"country"`
	Latitude  float64 `expr:"latitude"`
	Longitude float64 `expr:"longitude"`
}

// NewExprEnv creates an expression environment from a transaction event.
// It flattens common metadata fields for convenient access in expressions.
func NewExprEnv(event models.TransactionEvent) ExprEnv {
	env := ExprEnv{
		Amount:      event.Amount,
		Currency:    event.Currency,
		Type:        event.Type,
		FromAccount: event.FromAccount,
		ToAccount:   event.ToAccount,
		ExternalID:  event.ExternalID,
		Metadata:    event.Metadata,
	}

	// Extract common metadata fields for convenience
	if event.Metadata != nil {
		if region, ok := event.Metadata["region"].(string); ok {
			env.Region = region
		}
		if country, ok := event.Metadata["country"].(string); ok {
			env.Country = country
		}
		if lat, ok := toFloat64(event.Metadata["latitude"]); ok {
			env.Latitude = lat
		}
		if lon, ok := toFloat64(event.Metadata["longitude"]); ok {
			env.Longitude = lon
		}
	}

	return env
}

// EvaluateExpression compiles and evaluates an expression against a transaction event.
// Returns true if the expression evaluates to true, false otherwise.
// Uses expr-lang/expr for safe, type-checked expression evaluation.
func EvaluateExpression(expression string, event models.TransactionEvent) bool {
	if expression == "" {
		return false
	}

	// Build expression environment from transaction event
	env := NewExprEnv(event)

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
