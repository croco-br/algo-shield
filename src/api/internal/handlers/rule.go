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

type RuleHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewRuleHandler(db *pgxpool.Pool, redis *redis.Client) *RuleHandler {
	return &RuleHandler{
		db:    db,
		redis: redis,
	}
}

func (h *RuleHandler) CreateRule(c *fiber.Ctx) error {
	var rule models.Rule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	rule.ID = uuid.New()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	ctx := context.Background()
	conditionsJSON, _ := json.Marshal(rule.Conditions)

	query := `
		INSERT INTO rules (id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := h.db.Exec(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Type, rule.Action,
		rule.Priority, rule.Enabled, conditionsJSON, rule.Score,
		rule.CreatedAt, rule.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create rule",
		})
	}

	// Invalidate cache for hot-reload
	h.redis.Del(ctx, "rules:cache")

	return c.Status(fiber.StatusCreated).JSON(rule)
}

func (h *RuleHandler) GetRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	ctx := context.Background()
	var rule models.Rule
	var conditionsJSON []byte

	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		WHERE id = $1
	`

	err = h.db.QueryRow(ctx, query, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Action,
		&rule.Priority, &rule.Enabled, &conditionsJSON, &rule.Score,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found",
		})
	}

	json.Unmarshal(conditionsJSON, &rule.Conditions)

	return c.JSON(rule)
}

func (h *RuleHandler) ListRules(c *fiber.Ctx) error {
	ctx := context.Background()

	query := `
		SELECT id, name, description, type, action, priority, enabled, conditions, score, created_at, updated_at
		FROM rules
		ORDER BY priority ASC
	`

	rows, err := h.db.Query(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rules",
		})
	}
	defer rows.Close()

	rules := make([]models.Rule, 0)
	for rows.Next() {
		var rule models.Rule
		var conditionsJSON []byte

		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.Type, &rule.Action,
			&rule.Priority, &rule.Enabled, &conditionsJSON, &rule.Score,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(conditionsJSON, &rule.Conditions)
		rules = append(rules, rule)
	}

	return c.JSON(fiber.Map{
		"rules": rules,
	})
}

func (h *RuleHandler) UpdateRule(c *fiber.Ctx) error {
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
	rule.UpdatedAt = time.Now()

	ctx := context.Background()
	conditionsJSON, _ := json.Marshal(rule.Conditions)

	query := `
		UPDATE rules
		SET name = $2, description = $3, type = $4, action = $5, 
		    priority = $6, enabled = $7, conditions = $8, score = $9, updated_at = $10
		WHERE id = $1
	`

	result, err := h.db.Exec(ctx, query,
		rule.ID, rule.Name, rule.Description, rule.Type, rule.Action,
		rule.Priority, rule.Enabled, conditionsJSON, rule.Score, rule.UpdatedAt,
	)

	if err != nil || result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found or update failed",
		})
	}

	// Invalidate cache for hot-reload
	h.redis.Del(ctx, "rules:cache")

	return c.JSON(rule)
}

func (h *RuleHandler) DeleteRule(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid rule ID",
		})
	}

	ctx := context.Background()
	query := `DELETE FROM rules WHERE id = $1`

	result, err := h.db.Exec(ctx, query, id)
	if err != nil || result.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Rule not found or delete failed",
		})
	}

	// Invalidate cache for hot-reload
	h.redis.Del(ctx, "rules:cache")

	return c.SendStatus(fiber.StatusNoContent)
}
