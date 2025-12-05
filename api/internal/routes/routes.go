package routes

import (
	"github.com/algo-shield/algo-shield/api/internal/handlers"
	"github.com/algo-shield/algo-shield/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client) {
	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Handlers
	healthHandler := handlers.NewHealthHandler(db, redis)
	transactionHandler := handlers.NewTransactionHandler(db, redis)
	ruleHandler := handlers.NewRuleHandler(db, redis)

	// Health routes
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// API v1
	v1 := app.Group("/api/v1")

	// Transaction routes
	transactions := v1.Group("/transactions")
	transactions.Post("/", transactionHandler.ProcessTransaction)
	transactions.Get("/", transactionHandler.ListTransactions)
	transactions.Get("/:id", transactionHandler.GetTransaction)

	// Rule routes
	rules := v1.Group("/rules")
	rules.Post("/", ruleHandler.CreateRule)
	rules.Get("/", ruleHandler.ListRules)
	rules.Get("/:id", ruleHandler.GetRule)
	rules.Put("/:id", ruleHandler.UpdateRule)
	rules.Delete("/:id", ruleHandler.DeleteRule)
}

