package auth

import (
	"context"
	"strings"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	"github.com/algo-shield/algo-shield/src/pkg/csrf"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// csrfStoreTimeout is a dedicated timeout for CSRF token Redis storage.
// This must NOT reuse the handler context because bcrypt in login/register
// can consume nearly all of the handler timeout, leaving no time for the
// Redis SET operation — causing the CSRF token to silently not be stored.
const csrfStoreTimeout = 5 * time.Second

type Handler struct {
	service     AuthService
	userService UserService
	redis       *redis.Client
}

func NewHandler(service AuthService, userService UserService, redis *redis.Client) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
		redis:       redis,
	}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid request body"))
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return apierrors.SendError(c, apierrors.ValidationError(err.Error()))
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	// Register user (handles password hashing and token generation)
	user, token, refreshToken, err := h.service.RegisterUser(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		// Check if it's an APIError
		if apiErr, ok := err.(*apierrors.APIError); ok {
			return apierrors.SendError(c, apiErr)
		}
		// Check if it's a conflict error (user already exists)
		if err.Error() == "user with this email already exists" {
			return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrConflict, "User with this email already exists"))
		}
		return apierrors.SendError(c, apierrors.InternalError("Failed to register user"))
	}

	// SECURITY: Generate CSRF token for the new user
	csrfToken, err := csrf.GenerateToken()
	if err != nil {
		return apierrors.SendError(c, apierrors.InternalError("Failed to generate CSRF token"))
	}

	// Store CSRF token in Redis with a dedicated context (not the handler ctx
	// which may be nearly expired after bcrypt)
	csrfCtx, csrfCancel := context.WithTimeout(context.Background(), csrfStoreTimeout)
	defer csrfCancel()
	if err := csrf.StoreToken(csrfCtx, h.redis, user.ID.String(), csrfToken); err != nil {
		c.Context().Logger().Printf("Failed to store CSRF token: %v", err)
	}

	return c.JSON(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          user,
		"csrf_token":    csrfToken,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid request body"))
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return apierrors.SendError(c, apierrors.ValidationError(err.Error()))
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	// Login user (handles password verification and token generation)
	user, token, refreshToken, err := h.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		// Check if it's an APIError
		if apiErr, ok := err.(*apierrors.APIError); ok {
			return apierrors.SendError(c, apiErr)
		}
		return apierrors.SendError(c, apierrors.InternalError("Login failed"))
	}

	// SECURITY: Generate CSRF token for the logged-in user
	csrfToken, err := csrf.GenerateToken()
	if err != nil {
		return apierrors.SendError(c, apierrors.InternalError("Failed to generate CSRF token"))
	}

	// Store CSRF token in Redis with a dedicated context (not the handler ctx
	// which may be nearly expired after bcrypt)
	csrfCtx, csrfCancel := context.WithTimeout(context.Background(), csrfStoreTimeout)
	defer csrfCancel()
	if err := csrf.StoreToken(csrfCtx, h.redis, user.ID.String(), csrfToken); err != nil {
		c.Context().Logger().Printf("Failed to store CSRF token: %v", err)
	}

	return c.JSON(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          user,
		"csrf_token":    csrfToken,
	})
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "User not found in context"))
	}

	// Get fresh user data from service (with roles and groups loaded)
	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()
	freshUser, err := h.userService.GetUserByID(ctx, user.ID)
	if err != nil {
		return apierrors.SendError(c, apierrors.NotFound("User"))
	}

	return c.JSON(freshUser)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "Authorization header required"))
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "Invalid authorization header format"))
	}

	token := parts[1]

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	// Get user from context to delete CSRF token
	user, ok := c.Locals("user").(*models.User)
	if ok && user != nil {
		// SECURITY: Delete CSRF token from Redis
		if err := csrf.DeleteToken(ctx, h.redis, user.ID.String()); err != nil {
			c.Context().Logger().Printf("Failed to delete CSRF token on logout: %v", err)
		}
	}

	// Revoke the token
	if err := h.service.LogoutUser(ctx, token); err != nil {
		// Don't fail logout even if revocation fails
		// Log the error for debugging
		c.Context().Logger().Printf("Failed to revoke token on logout: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *Handler) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	return h.service.ValidateToken(ctx, tokenString)
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid request body"))
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return apierrors.SendError(c, apierrors.ValidationError(err.Error()))
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	// Refresh token (returns new access + refresh token via rotation)
	user, newAccessToken, newRefreshToken, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		// Check if it's an APIError
		if apiErr, ok := err.(*apierrors.APIError); ok {
			return apierrors.SendError(c, apiErr)
		}
		return apierrors.SendError(c, apierrors.InternalError("Failed to refresh token"))
	}

	// SECURITY: Generate new CSRF token
	csrfToken, err := csrf.GenerateToken()
	if err != nil {
		return apierrors.SendError(c, apierrors.InternalError("Failed to generate CSRF token"))
	}

	// Store CSRF token in Redis with a dedicated context
	csrfCtx, csrfCancel := context.WithTimeout(context.Background(), csrfStoreTimeout)
	defer csrfCancel()
	if err := csrf.StoreToken(csrfCtx, h.redis, user.ID.String(), csrfToken); err != nil {
		c.Context().Logger().Printf("Failed to store CSRF token: %v", err)
	}

	return c.JSON(fiber.Map{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
		"csrf_token":    csrfToken,
	})
}
