package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TransactionHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewTransactionHandler(db *pgxpool.Pool, redis *redis.Client) *TransactionHandler {
	return &TransactionHandler{
		db:    db,
		redis: redis,
	}
}

func (h *TransactionHandler) ProcessTransaction(c *fiber.Ctx) error {
	var event models.TransactionEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	startTime := time.Now()

	// Publish to Redis for worker processing
	ctx := context.Background()
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process transaction",
		})
	}

	// Push to Redis queue
	if err := h.redis.LPush(ctx, "transaction:queue", eventJSON).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to queue transaction",
		})
	}

	processingTime := time.Since(startTime).Milliseconds()

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":          "queued",
		"external_id":     event.ExternalID,
		"processing_time": processingTime,
		"message":         "Transaction queued for processing",
	})
}

func (h *TransactionHandler) GetTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	ctx := context.Background()
	var transaction models.Transaction

	query := `
		SELECT id, external_id, amount, currency, from_account, to_account, 
		       type, status, risk_score, risk_level, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		WHERE id = $1
	`

	err = h.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.ExternalID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.Type,
		&transaction.Status,
		&transaction.RiskScore,
		&transaction.RiskLevel,
		&transaction.ProcessingTime,
		&transaction.MatchedRules,
		&transaction.Metadata,
		&transaction.CreatedAt,
		&transaction.ProcessedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(transaction)
}

func (h *TransactionHandler) ListTransactions(c *fiber.Ctx) error {
	ctx := context.Background()

	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	query := `
		SELECT id, external_id, amount, currency, from_account, to_account, 
		       type, status, risk_score, risk_level, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := h.db.Query(ctx, query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.ExternalID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.FromAccount,
			&transaction.ToAccount,
			&transaction.Type,
			&transaction.Status,
			&transaction.RiskScore,
			&transaction.RiskLevel,
			&transaction.ProcessingTime,
			&transaction.MatchedRules,
			&transaction.Metadata,
			&transaction.CreatedAt,
			&transaction.ProcessedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, transaction)
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}
