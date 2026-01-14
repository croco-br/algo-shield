package transactions

import (
	"context"
	"fmt"
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()

	// Enqueue transaction using Asynq
	taskInfo, err := h.service.ProcessTransaction(ctx, event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to enqueue transaction",
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
		"job_id":      taskInfo.ID,
		"queue":       taskInfo.Queue,
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
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
	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
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
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid start_date format. Expected RFC3339 (e.g., 2006-01-02T15:04:05Z), got: %s", startDate),
			})
		}
		filter.StartDate = &t
	}
	if endDate := c.Query("end_date"); endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid end_date format. Expected RFC3339 (e.g., 2006-01-02T15:04:05Z), got: %s", endDate),
			})
		}
		filter.EndDate = &t
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	transaction, err := h.service.ApproveTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Transaction not found or not in pending/in_review status",
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.GetHandlerTimeout())
	defer cancel()
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	transaction, err := h.service.RejectTransaction(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Transaction not found or not in pending/in_review status",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reject transaction",
		})
	}

	return c.JSON(transaction)
}
