package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Validator instance with custom validations
	validate *validator.Validate

	// Email regex pattern - RFC 5322 compliant (simplified but robust)
	// Matches: user@domain.com, user.name@domain.co.uk, user+tag@example.com
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	registerCustomValidations()
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, formatError(e))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	return nil
}

// formatError formats a validation error into a user-friendly message
func formatError(e validator.FieldError) string {
	field := e.Field()
	if field == "" {
		field = e.StructField()
	}

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email", "email_regex":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "account":
		return fmt.Sprintf("%s must be a valid account identifier (alphanumeric, 1-100 chars)", field)
	case "currency":
		return fmt.Sprintf("%s must be a valid 3-letter currency code (ISO 4217)", field)
	case "transaction_type":
		return fmt.Sprintf("%s must be a valid transaction type", field)
	case "password":
		return fmt.Sprintf("%s must be between 8 and 128 characters", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// registerCustomValidations registers custom validation functions
func registerCustomValidations() {
	// Account validation: alphanumeric, 1-100 characters
	validate.RegisterValidation("account", func(fl validator.FieldLevel) bool {
		account := fl.Field().String()
		if len(account) < 1 || len(account) > 100 {
			return false
		}
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, account)
		return matched
	})

	// Currency validation: ISO 4217 (3 uppercase letters)
	validate.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		currency := fl.Field().String()
		matched, _ := regexp.MatchString(`^[A-Z]{3}$`, currency)
		return matched
	})

	// Transaction type validation
	validate.RegisterValidation("transaction_type", func(fl validator.FieldLevel) bool {
		transactionType := fl.Field().String()
		validTypes := []string{"transfer", "payment", "withdrawal", "deposit", "refund", "chargeback"}
		for _, validType := range validTypes {
			if transactionType == validType {
				return true
			}
		}
		return false
	})

	// Password validation: at least 8 chars, can be relaxed for development
	// In production, enforce: uppercase, lowercase, number, special char
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		// Minimum 8 characters required
		if len(password) < 8 {
			return false
		}
		// Maximum 128 characters to prevent DoS
		if len(password) > 128 {
			return false
		}
		// For now, just check length. Can be enhanced for production with:
		// hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
		// hasLower, _ := regexp.MatchString(`[a-z]`, password)
		// hasNumber, _ := regexp.MatchString(`[0-9]`, password)
		// hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
		// return hasUpper && hasLower && hasNumber && hasSpecial
		return true
	})

	// Email validation with regex: RFC 5322 compliant (simplified but robust)
	// Validates format: local-part@domain.tld
	// Supports: user@domain.com, user.name@domain.co.uk, user+tag@example.com
	validate.RegisterValidation("email_regex", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		// Check length constraints (RFC 5321)
		if len(email) < 3 || len(email) > 254 {
			return false
		}
		// Check basic structure: must contain @
		if !strings.Contains(email, "@") {
			return false
		}
		// Split email into local and domain parts
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		localPart := parts[0]
		domain := parts[1]

		// Validate local part (before @)
		if len(localPart) < 1 || len(localPart) > 64 {
			return false
		}
		// Local part cannot start or end with dot
		if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
			return false
		}
		// Local part cannot have consecutive dots
		if strings.Contains(localPart, "..") {
			return false
		}

		// Validate domain part (after @)
		if len(domain) < 1 || len(domain) > 253 {
			return false
		}
		// Domain cannot start or end with dot or hyphen
		if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
			strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
			return false
		}
		// Domain must contain at least one dot (for TLD)
		if !strings.Contains(domain, ".") {
			return false
		}
		// Domain cannot have consecutive dots
		if strings.Contains(domain, "..") {
			return false
		}

		// Use regex for final validation
		return emailRegex.MatchString(email)
	})
}

// ValidateLimit validates pagination limit parameter
func ValidateLimit(limit int) error {
	if limit < 1 {
		return fmt.Errorf("limit must be at least 1")
	}
	if limit > 1000 {
		return fmt.Errorf("limit must be at most 1000")
	}
	return nil
}

// ValidateOffset validates pagination offset parameter
func ValidateOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	return nil
}
