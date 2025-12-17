package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// DatabaseHealthChecker defines interface for database health checks
// Follows Dependency Inversion Principle
type DatabaseHealthChecker interface {
	Ping(ctx context.Context) error
}

// RedisHealthChecker defines interface for Redis health checks
// Follows Dependency Inversion Principle
type RedisHealthChecker interface {
	Ping(ctx context.Context) *redis.StatusCmd
}

type Handler struct {
	db    DatabaseHealthChecker
	redis RedisHealthChecker
}

// NewHandler creates a new health handler with dependency injection
// Follows Dependency Inversion Principle - receives interfaces, not concrete types
func NewHandler(db DatabaseHealthChecker, redis RedisHealthChecker) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	health := fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
	}

	// Check PostgreSQL
	if err := h.db.Ping(ctx); err != nil {
		health["postgres"] = "unhealthy"
		health["status"] = "degraded"
	} else {
		health["postgres"] = "healthy"
	}

	// Check Redis
	if err := h.redis.Ping(ctx).Err(); err != nil {
		health["redis"] = "unhealthy"
		health["status"] = "degraded"
	} else {
		health["redis"] = "healthy"
	}

	statusCode := fiber.StatusOK
	if health["status"] == "degraded" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(health)
}

func (h *Handler) Ready(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ready",
	})
}
