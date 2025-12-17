package permissions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// AssignRole assigns a role to a user
func (h *Handler) AssignRole(c *fiber.Ctx) error {
	userIDParam := c.Params("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.service.AssignRole(ctx, userID, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role assigned successfully",
	})
}

// RemoveRole removes a role from a user
func (h *Handler) RemoveRole(c *fiber.Ctx) error {
	userIDParam := c.Params("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	roleIDParam := c.Params("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.service.RemoveRole(ctx, userID, roleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role removed successfully",
	})
}

// ListRoles returns all available roles
func (h *Handler) ListRoles(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	roles, err := h.service.ListRoles(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}

	return c.JSON(fiber.Map{
		"roles": roles,
	})
}

// ListGroups returns all available groups
func (h *Handler) ListGroups(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	groups, err := h.service.ListGroups(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch groups",
		})
	}

	return c.JSON(fiber.Map{
		"groups": groups,
	})
}

// UpdateUserActive updates user active status
func (h *Handler) UpdateUserActive(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req UpdateUserActiveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.service.UpdateUserActive(ctx, userID, req.Active); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
