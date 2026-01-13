package dashboard

import (
	"context"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetMetrics returns aggregated dashboard metrics
// GET /api/v1/dashboard/metrics
func (h *Handler) GetMetrics(c *fiber.Ctx) error {
	start := time.Now()

	ctx, cancel := context.WithTimeout(c.Context(), internal.DEFAULT_TIMEOUT)
	defer cancel()
	ctx = shared.WithSyntheticMode(ctx, middleware.IsSyntheticMode(c))

	metrics, err := h.service.GetMetrics(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch dashboard metrics",
		})
	}

	// Add response time for performance monitoring
	responseTime := time.Since(start).Milliseconds()

	return c.JSON(fiber.Map{
		"data":             metrics,
		"response_time_ms": responseTime,
	})
}
