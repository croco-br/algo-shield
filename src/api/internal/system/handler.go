package system

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/gofiber/fiber/v2"
)

// Handler handles system configuration HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new system handler with dependency injection
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetSyntheticMode handles GET /api/v1/system/mode
func (h *Handler) GetSyntheticMode(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	response, err := h.service.GetSyntheticMode(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get synthetic mode",
		})
	}

	return c.JSON(response)
}

// SetSyntheticMode handles PUT /api/v1/system/mode
func (h *Handler) SetSyntheticMode(c *fiber.Ctx) error {
	var req UpdateSyntheticModeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	response, err := h.service.SetSyntheticMode(ctx, req.Enabled)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update synthetic mode",
		})
	}

	return c.JSON(response)
}
