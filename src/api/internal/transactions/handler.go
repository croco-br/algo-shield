package transactions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/api/internal"
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

func (h *Handler) ProcessTransaction(c *fiber.Ctx) error {
	var event models.TransactionEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.service.ProcessTransaction(ctx, event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to queue transaction",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":      "queued",
		"external_id": event.ExternalID,
		"message":     "Transaction queued for processing",
	})
}

func (h *Handler) GetTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	transaction, err := h.service.GetTransaction(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}

func (h *Handler) ListTransactions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	transactions, err := h.service.ListTransactions(ctx, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}
