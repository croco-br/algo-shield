package errors

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_NewAPIError_WhenCalled_ThenCreatesError(t *testing.T) {
	err := NewAPIError(ErrNotFound, "Resource not found")

	assert.Equal(t, ErrNotFound, err.Code)
	assert.Equal(t, "Resource not found", err.Message)
	assert.Empty(t, err.Details)
}

func Test_NewAPIErrorWithDetails_WhenCalled_ThenCreatesErrorWithDetails(t *testing.T) {
	err := NewAPIErrorWithDetails(ErrBadRequest, "Invalid input", "Field 'name' is required")

	assert.Equal(t, ErrBadRequest, err.Code)
	assert.Equal(t, "Invalid input", err.Message)
	assert.Equal(t, "Field 'name' is required", err.Details)
}

func Test_APIError_Error_WhenCalled_ThenReturnsMessage(t *testing.T) {
	err := NewAPIError(ErrNotFound, "Resource not found")

	result := err.Error()

	assert.Equal(t, "Resource not found", result)
}

func Test_GetHTTPStatus_WhenAuthErrors_ThenReturns401(t *testing.T) {
	authCodes := []ErrorCode{
		ErrInvalidCredentials,
		ErrTokenExpired,
		ErrTokenRevoked,
		ErrTokenInvalid,
		ErrUnauthorized,
	}

	for _, code := range authCodes {
		status := GetHTTPStatus(code)

		assert.Equal(t, fiber.StatusUnauthorized, status, "code %s should return 401", code)
	}
}

func Test_GetHTTPStatus_WhenPermissionErrors_ThenReturns403(t *testing.T) {
	forbiddenCodes := []ErrorCode{
		ErrUserInactive,
		ErrInsufficientPermissions,
		ErrCannotDeactivateSelf,
		ErrCannotDeactivateLastAdmin,
		ErrCannotModifyProtectedUser,
		ErrForbidden,
	}

	for _, code := range forbiddenCodes {
		status := GetHTTPStatus(code)

		assert.Equal(t, fiber.StatusForbidden, status, "code %s should return 403", code)
	}
}

func Test_GetHTTPStatus_WhenNotFound_ThenReturns404(t *testing.T) {
	status := GetHTTPStatus(ErrNotFound)

	assert.Equal(t, fiber.StatusNotFound, status)
}

func Test_GetHTTPStatus_WhenConflict_ThenReturns409(t *testing.T) {
	status := GetHTTPStatus(ErrConflict)

	assert.Equal(t, fiber.StatusConflict, status)
}

func Test_GetHTTPStatus_WhenRateLimit_ThenReturns429(t *testing.T) {
	status := GetHTTPStatus(ErrRateLimitExceeded)

	assert.Equal(t, fiber.StatusTooManyRequests, status)
}

func Test_GetHTTPStatus_WhenBadRequest_ThenReturns400(t *testing.T) {
	badRequestCodes := []ErrorCode{ErrBadRequest, ErrValidationError}

	for _, code := range badRequestCodes {
		status := GetHTTPStatus(code)

		assert.Equal(t, fiber.StatusBadRequest, status, "code %s should return 400", code)
	}
}

func Test_GetHTTPStatus_WhenInternalError_ThenReturns500(t *testing.T) {
	status := GetHTTPStatus(ErrInternalError)

	assert.Equal(t, fiber.StatusInternalServerError, status)
}

func Test_GetHTTPStatus_WhenUnknownCode_ThenReturns500(t *testing.T) {
	status := GetHTTPStatus(ErrorCode("UNKNOWN_CODE"))

	assert.Equal(t, fiber.StatusInternalServerError, status)
}

func Test_SendError_WhenCalled_ThenSendsJSONResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return SendError(c, NewAPIError(ErrNotFound, "Resource not found"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"code":"NOT_FOUND"`)
	assert.Contains(t, string(body), `"message":"Resource not found"`)
}

func Test_InvalidCredentials_WhenCalled_ThenReturnsError(t *testing.T) {
	err := InvalidCredentials()

	assert.Equal(t, ErrInvalidCredentials, err.Code)
	assert.Equal(t, "Invalid email or password", err.Message)
}

func Test_UserInactive_WhenCalled_ThenReturnsError(t *testing.T) {
	err := UserInactive()

	assert.Equal(t, ErrUserInactive, err.Code)
	assert.Equal(t, "User account is inactive", err.Message)
}

func Test_TokenRevoked_WhenCalled_ThenReturnsError(t *testing.T) {
	err := TokenRevoked()

	assert.Equal(t, ErrTokenRevoked, err.Code)
	assert.Equal(t, "Token has been revoked", err.Message)
}

func Test_TokenExpired_WhenCalled_ThenReturnsError(t *testing.T) {
	err := TokenExpired()

	assert.Equal(t, ErrTokenExpired, err.Code)
	assert.Equal(t, "Token has expired", err.Message)
}

func Test_TokenInvalid_WhenCalled_ThenReturnsError(t *testing.T) {
	err := TokenInvalid()

	assert.Equal(t, ErrTokenInvalid, err.Code)
	assert.Equal(t, "Invalid or malformed token", err.Message)
}

func Test_InsufficientPermissions_WhenCalled_ThenReturnsError(t *testing.T) {
	err := InsufficientPermissions()

	assert.Equal(t, ErrInsufficientPermissions, err.Code)
	assert.Contains(t, err.Message, "permission")
}

func Test_CannotDeactivateSelf_WhenCalled_ThenReturnsError(t *testing.T) {
	err := CannotDeactivateSelf()

	assert.Equal(t, ErrCannotDeactivateSelf, err.Code)
	assert.Contains(t, err.Message, "cannot deactivate")
}

func Test_CannotDeactivateLastAdmin_WhenCalled_ThenReturnsError(t *testing.T) {
	err := CannotDeactivateLastAdmin()

	assert.Equal(t, ErrCannotDeactivateLastAdmin, err.Code)
	assert.Contains(t, err.Message, "last active administrator")
}

func Test_NotFound_WhenCalled_ThenReturnsError(t *testing.T) {
	err := NotFound("User")

	assert.Equal(t, ErrNotFound, err.Code)
	assert.Equal(t, "User not found", err.Message)
}

func Test_InternalError_WhenCalled_ThenReturnsError(t *testing.T) {
	err := InternalError("Database connection failed")

	assert.Equal(t, ErrInternalError, err.Code)
	assert.Equal(t, "Database connection failed", err.Message)
}

func Test_BadRequest_WhenCalled_ThenReturnsError(t *testing.T) {
	err := BadRequest("Invalid parameter")

	assert.Equal(t, ErrBadRequest, err.Code)
	assert.Equal(t, "Invalid parameter", err.Message)
}

func Test_ValidationError_WhenCalled_ThenReturnsError(t *testing.T) {
	err := ValidationError("Field is required")

	assert.Equal(t, ErrValidationError, err.Code)
	assert.Equal(t, "Field is required", err.Message)
}
