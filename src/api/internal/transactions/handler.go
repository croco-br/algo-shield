package transactions

import (
	"context"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/middleware"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service Service
}

// NewHandler creates a new transaction handler with dependency injection
// Follows Dependency Inversion Principle - receives interface, not concrete type
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ProcessTransaction(c *fiber.Ctx) error {
	var event models.Event
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation: event must be a non-empty JSON object
	if len(event) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event must be a non-empty JSON object",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.service.ProcessTransaction(ctx, event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to queue transaction",
		})
	}

	// Try to extract external_id for response (if present)
	externalID := "unknown"
	if id, ok := event["external_id"].(string); ok {
		externalID = id
	} else if id, ok := event["id"].(string); ok {
		externalID = id
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":      "queued",
		"external_id": externalID,
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
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))
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
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate pagination parameters
	if err := validation.ValidateLimit(limit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := validation.ValidateOffset(offset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Build filter from query params
	filter := TransactionFilter{}

	if status := c.Query("status"); status != "" {
		s := models.TransactionStatus(status)
		filter.Status = &s
	}
	if schemaID := c.Query("schema_id"); schemaID != "" {
		if id, err := uuid.Parse(schemaID); err == nil {
			filter.SchemaID = &id
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			filter.EndDate = &t
		}
	}
	if minAmount := c.QueryFloat("min_amount", 0); minAmount > 0 {
		filter.MinAmount = &minAmount
	}
	if maxAmount := c.QueryFloat("max_amount", 0); maxAmount > 0 {
		filter.MaxAmount = &maxAmount
	}

	// Always use filtered method to get total count
	transactions, total, err := h.service.ListTransactionsWithFilter(ctx, filter, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}
	return c.JSON(fiber.Map{
		"transactions": transactions,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *Handler) ApproveTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	transaction, err := h.service.ApproveTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Transaction not found or not in review status",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve transaction",
		})
	}

	return c.JSON(transaction)
}

func (h *Handler) RejectTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	transaction, err := h.service.RejectTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Transaction not found or not in review status",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reject transaction",
		})
	}

	return c.JSON(transaction)
}
