package routes

import (
	"github.com/algo-shield/algo-shield/src/api/internal/handlers"
	"github.com/algo-shield/algo-shield/src/api/internal/middleware"
	"github.com/algo-shield/algo-shield/src/api/internal/services"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config) {
	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Services
	userService := services.NewUserService(db)

	// Handlers
	healthHandler := handlers.NewHealthHandler(db, redis)
	transactionHandler := handlers.NewTransactionHandler(db, redis)
	ruleHandler := handlers.NewRuleHandler(db, redis)
	authHandler := handlers.NewAuthHandler(userService, cfg)
	permissionsHandler := handlers.NewPermissionsHandler(userService)

	// Health routes (public)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// Auth routes (public)
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// API v1 (protected)
	v1 := app.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(authHandler))

	// Current user
	v1.Get("/auth/me", authHandler.GetCurrentUser)
	v1.Post("/auth/logout", authHandler.Logout)

	// Transaction routes (protected)
	transactions := v1.Group("/transactions")
	transactions.Post("/", transactionHandler.ProcessTransaction)
	transactions.Get("/", transactionHandler.ListTransactions)
	transactions.Get("/:id", transactionHandler.GetTransaction)

	// Rule routes (protected)
	rules := v1.Group("/rules")
	rules.Get("/", ruleHandler.ListRules)
	rules.Get("/:id", ruleHandler.GetRule)

	// Rule modification requires rule_editor or admin role
	rulesProtected := rules.Group("", middleware.RequireAnyRole("admin", "rule_editor"))
	rulesProtected.Post("/", ruleHandler.CreateRule)
	rulesProtected.Put("/:id", ruleHandler.UpdateRule)
	rulesProtected.Delete("/:id", ruleHandler.DeleteRule)

	// Permissions management (admin only)
	permissions := v1.Group("/permissions", middleware.RequireRole("admin"))
	permissions.Get("/users", permissionsHandler.ListUsers)
	permissions.Get("/users/:id", permissionsHandler.GetUser)
	permissions.Put("/users/:id/active", permissionsHandler.UpdateUserActive)
	permissions.Post("/users/:userId/roles", permissionsHandler.AssignRole)
	permissions.Delete("/users/:userId/roles/:roleId", permissionsHandler.RemoveRole)
	permissions.Get("/roles", permissionsHandler.ListRoles)
	permissions.Get("/groups", permissionsHandler.ListGroups)
}
