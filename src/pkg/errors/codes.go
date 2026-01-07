package errors

import "github.com/gofiber/fiber/v2"

// ErrorCode represents a language-agnostic error code for i18n support
type ErrorCode string

const (
	// Authentication errors
	ErrInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrUserInactive       ErrorCode = "USER_INACTIVE"
	ErrTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	ErrTokenRevoked       ErrorCode = "TOKEN_REVOKED"
	ErrTokenInvalid       ErrorCode = "TOKEN_INVALID"

	// Permission errors
	ErrInsufficientPermissions   ErrorCode = "INSUFFICIENT_PERMISSIONS"
	ErrCannotDeactivateSelf      ErrorCode = "CANNOT_DEACTIVATE_SELF"
	ErrCannotDeactivateLastAdmin ErrorCode = "CANNOT_DEACTIVATE_LAST_ADMIN"
	ErrCannotModifyProtectedUser ErrorCode = "CANNOT_MODIFY_PROTECTED_USER"

	// Rate limiting
	ErrRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Generic errors
	ErrNotFound        ErrorCode = "NOT_FOUND"
	ErrInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrBadRequest      ErrorCode = "BAD_REQUEST"
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrForbidden       ErrorCode = "FORBIDDEN"
	ErrConflict        ErrorCode = "CONFLICT"
	ErrValidationError ErrorCode = "VALIDATION_ERROR"
)

// APIError represents a structured API error with code and message
type APIError struct {
	Code    ErrorCode `json:"code"`              // Language-agnostic error code
	Message string    `json:"message"`           // English fallback message
	Details string    `json:"details,omitempty"` // Optional details (only in dev mode)
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// NewAPIErrorWithDetails creates a new APIError with details
func NewAPIErrorWithDetails(code ErrorCode, message, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// GetHTTPStatus returns the appropriate HTTP status code for an error code
func GetHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrInvalidCredentials, ErrTokenExpired, ErrTokenRevoked, ErrTokenInvalid, ErrUnauthorized:
		return fiber.StatusUnauthorized
	case ErrUserInactive, ErrInsufficientPermissions, ErrCannotDeactivateSelf,
		ErrCannotDeactivateLastAdmin, ErrCannotModifyProtectedUser, ErrForbidden:
		return fiber.StatusForbidden
	case ErrNotFound:
		return fiber.StatusNotFound
	case ErrConflict:
		return fiber.StatusConflict
	case ErrRateLimitExceeded:
		return fiber.StatusTooManyRequests
	case ErrBadRequest, ErrValidationError:
		return fiber.StatusBadRequest
	case ErrInternalError:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

// SendError sends an APIError as a JSON response
func SendError(c *fiber.Ctx, err *APIError) error {
	status := GetHTTPStatus(err.Code)
	return c.Status(status).JSON(err)
}

// Common error constructors for convenience
func InvalidCredentials() *APIError {
	return NewAPIError(ErrInvalidCredentials, "Invalid email or password")
}

func UserInactive() *APIError {
	return NewAPIError(ErrUserInactive, "User account is inactive")
}

func TokenRevoked() *APIError {
	return NewAPIError(ErrTokenRevoked, "Token has been revoked")
}

func TokenExpired() *APIError {
	return NewAPIError(ErrTokenExpired, "Token has expired")
}

func TokenInvalid() *APIError {
	return NewAPIError(ErrTokenInvalid, "Invalid or malformed token")
}

func InsufficientPermissions() *APIError {
	return NewAPIError(ErrInsufficientPermissions, "You don't have permission to perform this action")
}

func CannotDeactivateSelf() *APIError {
	return NewAPIError(ErrCannotDeactivateSelf, "You cannot deactivate your own account")
}

func CannotDeactivateLastAdmin() *APIError {
	return NewAPIError(ErrCannotDeactivateLastAdmin, "Cannot deactivate the last active administrator")
}

func NotFound(resource string) *APIError {
	return NewAPIError(ErrNotFound, resource+" not found")
}

func InternalError(message string) *APIError {
	return NewAPIError(ErrInternalError, message)
}

func BadRequest(message string) *APIError {
	return NewAPIError(ErrBadRequest, message)
}

func ValidationError(message string) *APIError {
	return NewAPIError(ErrValidationError, message)
}
