package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	RequiredField string `json:"required_field" validate:"required"`
	EmailField    string `json:"email_field" validate:"email_regex"`
	MinField      string `json:"min_field" validate:"min=5"`
	MaxField      string `json:"max_field" validate:"max=10"`
}

type AccountStruct struct {
	Account string `json:"account" validate:"account"`
}

type CurrencyStruct struct {
	Currency string `json:"currency" validate:"currency"`
}

type TransactionTypeStruct struct {
	Type string `json:"type" validate:"transaction_type"`
}

type PasswordStruct struct {
	Password string `json:"password" validate:"password"`
}

func Test_ValidateStruct_WhenValidStruct_ThenReturnsNil(t *testing.T) {
	valid := TestStruct{
		RequiredField: "value",
		EmailField:    "test@example.com",
		MinField:      "12345",
		MaxField:      "1234567890",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenRequiredFieldMissing_ThenReturnsError(t *testing.T) {
	invalid := TestStruct{
		EmailField: "test@example.com",
		MinField:   "12345",
		MaxField:   "1234567890",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required_field is required")
}

func Test_ValidateStruct_WhenEmailInvalid_ThenReturnsError(t *testing.T) {
	invalid := TestStruct{
		RequiredField: "value",
		EmailField:    "invalid-email",
		MinField:      "12345",
		MaxField:      "1234567890",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email_field")
}

func Test_ValidateStruct_WhenMinLengthViolated_ThenReturnsError(t *testing.T) {
	invalid := TestStruct{
		RequiredField: "value",
		EmailField:    "test@example.com",
		MinField:      "123",
		MaxField:      "1234567890",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "min_field")
	assert.Contains(t, err.Error(), "at least")
}

func Test_ValidateStruct_WhenMaxLengthViolated_ThenReturnsError(t *testing.T) {
	invalid := TestStruct{
		RequiredField: "value",
		EmailField:    "test@example.com",
		MinField:      "12345",
		MaxField:      "12345678901",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max_field")
	assert.Contains(t, err.Error(), "at most")
}

func Test_ValidateStruct_WhenValidAccount_ThenReturnsNil(t *testing.T) {
	valid := AccountStruct{
		Account: "user123",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenAccountWithHyphen_ThenReturnsNil(t *testing.T) {
	valid := AccountStruct{
		Account: "user-123",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenAccountWithUnderscore_ThenReturnsNil(t *testing.T) {
	valid := AccountStruct{
		Account: "user_123",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenAccountTooLong_ThenReturnsError(t *testing.T) {
	invalid := AccountStruct{
		Account: strings.Repeat("a", 101),
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account")
}

func Test_ValidateStruct_WhenAccountEmpty_ThenReturnsError(t *testing.T) {
	invalid := AccountStruct{
		Account: "",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
}

func Test_ValidateStruct_WhenAccountInvalidChars_ThenReturnsError(t *testing.T) {
	invalid := AccountStruct{
		Account: "user@123",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
}

func Test_ValidateStruct_WhenValidCurrency_ThenReturnsNil(t *testing.T) {
	valid := CurrencyStruct{
		Currency: "USD",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenCurrencyLowercase_ThenReturnsError(t *testing.T) {
	invalid := CurrencyStruct{
		Currency: "usd",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "currency")
}

func Test_ValidateStruct_WhenCurrencyTooShort_ThenReturnsError(t *testing.T) {
	invalid := CurrencyStruct{
		Currency: "US",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
}

func Test_ValidateStruct_WhenCurrencyTooLong_ThenReturnsError(t *testing.T) {
	invalid := CurrencyStruct{
		Currency: "USDD",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
}

func Test_ValidateStruct_WhenValidTransactionType_ThenReturnsNil(t *testing.T) {
	validTypes := []string{"transfer", "payment", "withdrawal", "deposit", "refund", "chargeback"}

	for _, validType := range validTypes {
		valid := TransactionTypeStruct{
			Type: validType,
		}

		err := ValidateStruct(valid)

		assert.NoError(t, err, "type %s should be valid", validType)
	}
}

func Test_ValidateStruct_WhenInvalidTransactionType_ThenReturnsError(t *testing.T) {
	invalid := TransactionTypeStruct{
		Type: "invalid_type",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type")
}

func Test_ValidateStruct_WhenValidPassword_ThenReturnsNil(t *testing.T) {
	valid := PasswordStruct{
		Password: "12345678",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenPasswordTooShort_ThenReturnsError(t *testing.T) {
	invalid := PasswordStruct{
		Password: "1234567",
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password")
}

func Test_ValidateStruct_WhenPasswordTooLong_ThenReturnsError(t *testing.T) {
	invalid := PasswordStruct{
		Password: strings.Repeat("a", 129),
	}

	err := ValidateStruct(invalid)

	assert.Error(t, err)
}

func Test_ValidateStruct_WhenPasswordAtMinLength_ThenReturnsNil(t *testing.T) {
	valid := PasswordStruct{
		Password: "12345678",
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateStruct_WhenPasswordAtMaxLength_ThenReturnsNil(t *testing.T) {
	valid := PasswordStruct{
		Password: strings.Repeat("a", 128),
	}

	err := ValidateStruct(valid)

	assert.NoError(t, err)
}

func Test_ValidateLimit_WhenValidLimit_ThenReturnsNil(t *testing.T) {
	err := ValidateLimit(50)

	assert.NoError(t, err)
}

func Test_ValidateLimit_WhenLimitTooLow_ThenReturnsError(t *testing.T) {
	err := ValidateLimit(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 1")
}

func Test_ValidateLimit_WhenLimitTooHigh_ThenReturnsError(t *testing.T) {
	err := ValidateLimit(1001)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at most 1000")
}

func Test_ValidateLimit_WhenLimitAt1_ThenReturnsNil(t *testing.T) {
	err := ValidateLimit(1)

	assert.NoError(t, err)
}

func Test_ValidateLimit_WhenLimitAt1000_ThenReturnsNil(t *testing.T) {
	err := ValidateLimit(1000)

	assert.NoError(t, err)
}

func Test_ValidateOffset_WhenValidOffset_ThenReturnsNil(t *testing.T) {
	err := ValidateOffset(0)

	assert.NoError(t, err)
}

func Test_ValidateOffset_WhenPositiveOffset_ThenReturnsNil(t *testing.T) {
	err := ValidateOffset(100)

	assert.NoError(t, err)
}

func Test_ValidateOffset_WhenNegativeOffset_ThenReturnsError(t *testing.T) {
	err := ValidateOffset(-1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-negative")
}

func Test_NewValidationError_WhenCalled_ThenReturnsValidationError(t *testing.T) {
	err := NewValidationError("test error")

	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

func Test_IsValidationError_WhenValidationError_ThenReturnsTrue(t *testing.T) {
	err := NewValidationError("test error")

	result := IsValidationError(err)

	assert.True(t, result)
}

func Test_IsValidationError_WhenNotValidationError_ThenReturnsFalse(t *testing.T) {
	err := assert.AnError

	result := IsValidationError(err)

	assert.False(t, result)
}

func Test_ValidateStruct_WhenValidEmail_ThenReturnsNil(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"simple", "user@example.com"},
		{"with dot", "user.name@example.com"},
		{"with plus", "user+tag@example.com"},
		{"subdomain", "user@mail.example.com"},
		{"two letter TLD", "user@example.co"},
		{"three letter TLD", "user@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := TestStruct{
				RequiredField: "value",
				EmailField:    tt.email,
				MinField:      "12345",
				MaxField:      "1234567890",
			}

			err := ValidateStruct(valid)

			assert.NoError(t, err)
		})
	}
}

func Test_ValidateStruct_WhenInvalidEmail_ThenReturnsError(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"no at", "userexample.com"},
		{"no domain", "user@"},
		{"no local", "@example.com"},
		{"no TLD", "user@example"},
		{"double at", "user@@example.com"},
		{"starts with dot", ".user@example.com"},
		{"ends with dot", "user.@example.com"},
		{"consecutive dots", "user..name@example.com"},
		{"domain starts with dot", "user@.example.com"},
		{"domain ends with dot", "user@example.com."},
		{"too long", strings.Repeat("a", 255) + "@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invalid := TestStruct{
				RequiredField: "value",
				EmailField:    tt.email,
				MinField:      "12345",
				MaxField:      "1234567890",
			}

			err := ValidateStruct(invalid)

			assert.Error(t, err)
		})
	}
}

func Test_ValidationError_Error_WhenCalled_ThenReturnsMessage(t *testing.T) {
	err := &ValidationError{Message: "test message"}

	result := err.Error()

	assert.Equal(t, "test message", result)
}
