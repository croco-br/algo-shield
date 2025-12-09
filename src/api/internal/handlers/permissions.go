package handlers

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PermissionsHandler struct {
	userService *services.UserService
}

func NewPermissionsHandler(userService *services.UserService) *PermissionsHandler {
	return &PermissionsHandler{
		userService: userService,
	}
}

// ListUsers returns all users with their roles and groups
func (h *PermissionsHandler) ListUsers(c *fiber.Ctx) error {
	ctx := context.Background()
	users, err := h.userService.ListUsers(ctx)
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
func (h *PermissionsHandler) GetUser(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// AssignRole assigns a role to a user
func (h *PermissionsHandler) AssignRole(c *fiber.Ctx) error {
	userIDParam := c.Params("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	if err := h.userService.AssignRole(ctx, userID, req.RoleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role assigned successfully",
	})
}

// RemoveRole removes a role from a user
func (h *PermissionsHandler) RemoveRole(c *fiber.Ctx) error {
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

	ctx := context.Background()
	if err := h.userService.RemoveRole(ctx, userID, roleID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Role removed successfully",
	})
}

// ListRoles returns all available roles
func (h *PermissionsHandler) ListRoles(c *fiber.Ctx) error {
	ctx := context.Background()
	roles, err := h.userService.ListRoles(ctx)
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
func (h *PermissionsHandler) ListGroups(c *fiber.Ctx) error {
	ctx := context.Background()
	groups, err := h.userService.ListGroups(ctx)
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
func (h *PermissionsHandler) UpdateUserActive(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req struct {
		Active bool `json:"active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	if err := h.userService.UpdateUserActive(ctx, userID, req.Active); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}
