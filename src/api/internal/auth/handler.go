package auth

import (
	"context"
	"strings"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service     AuthService
	userService UserService
}

func NewHandler(service AuthService, userService UserService) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	// Register user (handles password hashing and token generation)
	user, token, err := h.service.RegisterUser(ctx, req.Email, req.Name, req.Password)
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

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	// Login user (handles password verification and token generation)
	user, token, err := h.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		// Check if it's an APIError
		if apiErr, ok := err.(*apierrors.APIError); ok {
			return apierrors.SendError(c, apiErr)
		}
		return apierrors.SendError(c, apierrors.InternalError("Login failed"))
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "User not found in context"))
	}

	// Get fresh user data from service (with roles and groups loaded)
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

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

func (h *Handler) ValidateToken(tokenString string) (*models.User, error) {
	return h.service.ValidateToken(tokenString)
}
