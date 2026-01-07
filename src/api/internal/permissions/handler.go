package permissions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service PermissionsService
}

func NewHandler(service PermissionsService) *Handler {
	return &Handler{
		service: service,
	}
}

// ListUsers returns all users with their roles and groups
func (h *Handler) ListUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	users, err := h.service.ListUsers(ctx)
	if err != nil {
		return apierrors.SendError(c, apierrors.InternalError("Failed to fetch users"))
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

// GetUser returns a specific user by ID
func (h *Handler) GetUser(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid user ID"))
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		return apierrors.SendError(c, apierrors.NotFound("User"))
	}

	return c.JSON(user)
}

// UpdateUserActive updates user active status with admin protection
func (h *Handler) UpdateUserActive(c *fiber.Ctx) error {
	// Get current user from context
	currentUser, ok := c.Locals("user").(*models.User)
	if !ok {
		return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "User not found in context"))
	}

	// Parse target user ID
	targetUserIDParam := c.Params("id")
	targetUserID, err := uuid.Parse(targetUserIDParam)
	if err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid user ID"))
	}

	var req UpdateUserActiveRequest
	if err := c.BodyParser(&req); err != nil {
		return apierrors.SendError(c, apierrors.BadRequest("Invalid request body"))
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return apierrors.SendError(c, apierrors.ValidationError(err.Error()))
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	// Call service with admin protection
	// req.Active is guaranteed to be non-nil after validation
	if err := h.service.UpdateUserActive(ctx, currentUser.ID, targetUserID, *req.Active); err != nil {
		// Check if it's an APIError
		if apiErr, ok := err.(*apierrors.APIError); ok {
			return apierrors.SendError(c, apiErr)
		}
		return apierrors.SendError(c, apierrors.InternalError("Failed to update user"))
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
