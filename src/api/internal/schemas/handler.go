package schemas

import (
	"context"
	"errors"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for schemas
type Handler struct {
	service *Service
}

// NewHandler creates a new schema handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateSchema handles POST /api/v1/schemas
func (h *Handler) CreateSchema(c *fiber.Ctx) error {
	var req CreateSchemaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(req.SampleJSON) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sample_json is required and must be a non-empty JSON object",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	schema, err := h.service.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, ErrSchemaNameExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A schema with this name already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create schema",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(schema)
}

// GetSchema handles GET /api/v1/schemas/:id
func (h *Handler) GetSchema(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schema ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	schema, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch schema",
		})
	}

	return c.JSON(schema)
}

// ListSchemas handles GET /api/v1/schemas
func (h *Handler) ListSchemas(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	schemas, err := h.service.List(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch schemas",
		})
	}

	return c.JSON(SchemaListResponse{Schemas: schemas})
}

// UpdateSchema handles PUT /api/v1/schemas/:id
func (h *Handler) UpdateSchema(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schema ID",
		})
	}

	var req UpdateSchemaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validation.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	schema, err := h.service.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}
		if errors.Is(err, ErrSchemaNameExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A schema with this name already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update schema",
		})
	}

	return c.JSON(schema)
}

// DeleteSchema handles DELETE /api/v1/schemas/:id
func (h *Handler) DeleteSchema(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schema ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	err = h.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}
		if errors.Is(err, ErrSchemaHasRules) {
			// Get the list of rules referencing this schema
			rules, _ := h.service.GetRulesReferencingSchema(ctx, id)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":             "Schema is referenced by rules and cannot be deleted",
				"referencing_rules": rules,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete schema",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ParseSchema handles POST /api/v1/schemas/:id/parse
// Re-extracts fields from the schema's sample JSON
func (h *Handler) ParseSchema(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schema ID",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()

	schema, err := h.service.ParseSampleJSON(ctx, id)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse schema",
		})
	}

	return c.JSON(schema)
}
