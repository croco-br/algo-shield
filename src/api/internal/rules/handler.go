package rules

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	service *Service
}

func NewHandler(db *pgxpool.Pool, redis *redis.Client) *Handler {
	service := NewService(db, redis)
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateRule(c *fiber.Ctx) error {
	var rule models.Rule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	if err := h.service.CreateRule(ctx, &rule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create rule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(rule)
}

func (h *Handler) GetRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	ctx := context.Background()
	rule, err := h.service.GetRule(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found",
		})
	}

	return c.JSON(rule)
}

func (h *Handler) ListRules(c *fiber.Ctx) error {
	ctx := context.Background()

	rules, err := h.service.ListRules(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rules",
		})
	}

	return c.JSON(fiber.Map{
		"rules": rules,
	})
}

func (h *Handler) UpdateRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	var rule models.Rule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	rule.ID = id

	ctx := context.Background()
	if err := h.service.UpdateRule(ctx, &rule); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found or update failed",
		})
	}

	return c.JSON(rule)
}

func (h *Handler) DeleteRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	ctx := context.Background()
	if err := h.service.DeleteRule(ctx, id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found or delete failed",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
