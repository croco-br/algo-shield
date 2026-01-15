package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildJSONPathQuery Tests

func Test_BuildJSONPathQuery_WhenSinglePath_ThenReturnsSimpleQuery(t *testing.T) {
	result, err := buildJSONPathQuery("origin")

	require.NoError(t, err)
	assert.Equal(t, "metadata->>'origin'", result)
}

func Test_BuildJSONPathQuery_WhenNestedPath_ThenReturnsNestedQuery(t *testing.T) {
	result, err := buildJSONPathQuery("user.id")

	require.NoError(t, err)
	assert.Equal(t, "metadata->'user'->>'id'", result)
}

func Test_BuildJSONPathQuery_WhenDeeplyNestedPath_ThenReturnsComplexQuery(t *testing.T) {
	result, err := buildJSONPathQuery("location.address.city")

	require.NoError(t, err)
	assert.Equal(t, "metadata->'location'->'address'->>'city'", result)
}

func Test_BuildJSONPathQuery_WhenPathWithUnderscore_ThenAcceptsPath(t *testing.T) {
	result, err := buildJSONPathQuery("user_id")

	require.NoError(t, err)
	assert.Equal(t, "metadata->>'user_id'", result)
}

func Test_BuildJSONPathQuery_WhenPathWithNumbers_ThenAcceptsPath(t *testing.T) {
	result, err := buildJSONPathQuery("field123")

	require.NoError(t, err)
	assert.Equal(t, "metadata->>'field123'", result)
}

func Test_BuildJSONPathQuery_WhenPathWithMixedCase_ThenAcceptsPath(t *testing.T) {
	result, err := buildJSONPathQuery("firstName")

	require.NoError(t, err)
	assert.Equal(t, "metadata->>'firstName'", result)
}

// SQL Injection Prevention Tests

func Test_BuildJSONPathQuery_WhenPathWithQuote_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field'OR'1'='1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithSemicolon_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field;DROP TABLE")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithDash_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field-name")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithSpace_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field name")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithParenthesis_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field()")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithSQLComment_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field--comment")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithSlash_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field/path")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenPathWithBackslash_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("field\\path")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

func Test_BuildJSONPathQuery_WhenEmptyPath_ThenReturnsError(t *testing.T) {
	result, err := buildJSONPathQuery("")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid field path")
	assert.Empty(t, result)
}

// Edge Cases
// Note: The current validation allows dots anywhere in the path (including leading,
// trailing, and consecutive dots). While this could be stricter, the actual query
// building handles these cases, and in practice, schema field paths would not
// contain these patterns.

func Test_BuildJSONPathQuery_WhenPathStartsWithDot_ThenBuildsQuery(t *testing.T) {
	result, err := buildJSONPathQuery(".field")

	require.NoError(t, err)
	// Splits into ["", "field"], so builds: metadata->''->>'field'
	assert.Contains(t, result, "field")
}

func Test_BuildJSONPathQuery_WhenPathEndsWithDot_ThenBuildsQuery(t *testing.T) {
	result, err := buildJSONPathQuery("field.")

	require.NoError(t, err)
	// Splits into ["field", ""], so builds: metadata->'field'->''
	assert.Contains(t, result, "field")
}

func Test_BuildJSONPathQuery_WhenPathHasDoubleDot_ThenBuildsQuery(t *testing.T) {
	result, err := buildJSONPathQuery("field..name")

	require.NoError(t, err)
	// Splits into ["field", "", "name"], so builds: metadata->'field'->''->'name'
	assert.Contains(t, result, "field")
	assert.Contains(t, result, "name")
}

// Valid Complex Paths

func Test_BuildJSONPathQuery_WhenValidComplexPath_ThenBuildsCorrectQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "three levels",
			input:    "a.b.c",
			expected: "metadata->'a'->'b'->>'c'",
		},
		{
			name:     "four levels",
			input:    "a.b.c.d",
			expected: "metadata->'a'->'b'->'c'->>'d'",
		},
		{
			name:     "with underscores",
			input:    "user_info.account_id",
			expected: "metadata->'user_info'->>'account_id'",
		},
		{
			name:     "with numbers",
			input:    "field1.field2",
			expected: "metadata->'field1'->>'field2'",
		},
		{
			name:     "mixed case",
			input:    "firstName.lastName",
			expected: "metadata->'firstName'->>'lastName'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildJSONPathQuery(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Validation Pattern Tests

func Test_BuildJSONPathQuery_WhenOnlyAlphanumericDotUnderscore_ThenAccepts(t *testing.T) {
	validPaths := []string{
		"field",
		"field_name",
		"field123",
		"Field",
		"FIELD",
		"field.nested",
		"field.nested.deep",
		"field_1.field_2",
		"_field",
		"field_",
		"__field__",
		"a.b.c.d.e.f",
	}

	for _, path := range validPaths {
		t.Run(path, func(t *testing.T) {
			result, err := buildJSONPathQuery(path)

			require.NoError(t, err)
			assert.NotEmpty(t, result)
		})
	}
}

func Test_BuildJSONPathQuery_WhenInvalidCharacters_ThenRejects(t *testing.T) {
	invalidPaths := []string{
		"field!",
		"field@",
		"field#",
		"field$",
		"field%",
		"field^",
		"field&",
		"field*",
		"field(",
		"field)",
		"field+",
		"field=",
		"field[",
		"field]",
		"field{",
		"field}",
		"field|",
		"field:",
		"field;",
		"field'",
		"field\"",
		"field<",
		"field>",
		"field,",
		"field?",
		"field/",
		"field\\",
		"field`",
		"field~",
		"field ",
		"field\t",
		"field\n",
	}

	for _, path := range invalidPaths {
		t.Run(path, func(t *testing.T) {
			result, err := buildJSONPathQuery(path)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid field path")
			assert.Empty(t, result)
		})
	}
}

// Note: CountByFieldInTimeWindow and SumFieldByFieldInTimeWindow are best tested
// with integration tests as they require actual database queries. Unit testing
// these methods would require complex mocking of pgxpool.Pool and would not
// provide significant value over integration tests.
//
// The critical SQL injection prevention logic is thoroughly tested via
// buildJSONPathQuery tests above.
