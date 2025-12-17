package rules

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	repo rules.Repository
}

func NewHandler(db *pgxpool.Pool, redis *redis.Client) *Handler {
	return &Handler{
		repo: rules.NewPostgresRepository(db, redis),
	}
}

func (h *Handler) CreateRule(c *fiber.Ctx) error {
	var rule models.Rule
	if err := c.BodyParser(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate rule
	if err := validateRule(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set timestamps
	now := time.Now()
	if rule.ID == uuid.Nil {
		rule.ID = uuid.New()
	}
	rule.CreatedAt = now
	rule.UpdatedAt = now

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.repo.CreateRule(ctx, &rule); err != nil {
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	rule, err := h.repo.GetRule(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rule not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rule",
		})
	}

	return c.JSON(rule)
}

func (h *Handler) ListRules(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	rules, err := h.repo.ListRules(ctx)
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
	rule.UpdatedAt = time.Now()

	// Validate rule
	if err := validateRule(&rule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.repo.UpdateRule(ctx, &rule); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rule not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update rule",
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

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	if err := h.repo.DeleteRule(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Rule not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete rule",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// validateRule validates a rule struct and returns an error if invalid
func validateRule(rule *models.Rule) error {
	// Validate name
	if rule.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(rule.Name) > 255 {
		return fmt.Errorf("name must be at most 255 characters")
	}

	// Validate type
	validTypes := []models.RuleType{
		models.RuleTypeAmount,
		models.RuleTypeVelocity,
		models.RuleTypeBlocklist,
		models.RuleTypePattern,
		models.RuleTypeGeography,
		models.RuleTypeCustom,
	}
	typeValid := false
	for _, validType := range validTypes {
		if rule.Type == validType {
			typeValid = true
			break
		}
	}
	if !typeValid {
		return fmt.Errorf("invalid rule type: %s. Valid types are: amount, velocity, blocklist, pattern, geography, custom", rule.Type)
	}

	// Validate action
	validActions := []models.RuleAction{
		models.ActionAllow,
		models.ActionBlock,
		models.ActionReview,
		models.ActionScore,
	}
	actionValid := false
	for _, validAction := range validActions {
		if rule.Action == validAction {
			actionValid = true
			break
		}
	}
	if !actionValid {
		return fmt.Errorf("invalid rule action: %s. Valid actions are: allow, block, review, score", rule.Action)
	}

	// Validate priority
	if rule.Priority < 0 {
		return fmt.Errorf("priority must be non-negative")
	}
	if rule.Priority > 1000 {
		return fmt.Errorf("priority must be at most 1000")
	}

	// Validate score
	if rule.Score < 0 {
		return fmt.Errorf("score must be non-negative")
	}
	if rule.Score > 100 {
		return fmt.Errorf("score must be at most 100")
	}

	// Validate conditions (must not be nil, but can be empty)
	if rule.Conditions == nil {
		return fmt.Errorf("conditions cannot be nil")
	}

	return nil
}
