package groups

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
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

// GetGroup returns a specific group by ID
func (h *Handler) GetGroup(c *fiber.Ctx) error {
	groupIDParam := c.Params("id")
	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	group, err := h.service.GetGroupByID(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}

	// Load roles for the group
	roles, err := h.service.LoadGroupRoles(ctx, groupID)
	if err == nil {
		group.Roles = roles
	}

	return c.JSON(group)
}
