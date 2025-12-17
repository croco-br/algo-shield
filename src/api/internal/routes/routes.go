package routes

import (
	"github.com/algo-shield/algo-shield/src/api/internal/auth"
	"github.com/algo-shield/algo-shield/src/api/internal/health"
	"github.com/algo-shield/algo-shield/src/api/internal/permissions"
	"github.com/algo-shield/algo-shield/src/api/internal/rules"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/middleware"
	"github.com/algo-shield/algo-shield/src/api/internal/transactions"
	"github.com/algo-shield/algo-shield/src/api/internal/user"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config) {
	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	// Initialize slices
	userService := user.NewService(db, cfg)
	authService := auth.NewService(db, cfg, userService)
	authHandler := auth.NewHandler(authService, userService)

	permissionsService := permissions.NewService(db)
	permissionsHandler := permissions.NewHandler(permissionsService)

	healthHandler := health.NewHandler(db, redis)
	transactionHandler := transactions.NewHandler(db, redis)
	ruleHandler := rules.NewHandler(db, redis)

	// Health routes (public)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// Auth routes (public)
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)

	// API v1 (protected)
	v1 := app.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(authHandler))

	// Current user
	v1.Get("/auth/me", authHandler.GetCurrentUser)
	v1.Post("/auth/logout", authHandler.Logout)

	// Transaction routes (protected)
	transactionsGroup := v1.Group("/transactions")
	transactionsGroup.Post("/", transactionHandler.ProcessTransaction)
	transactionsGroup.Get("/", transactionHandler.ListTransactions)
	transactionsGroup.Get("/:id", transactionHandler.GetTransaction)

	// Rule routes (protected)
	rulesGroup := v1.Group("/rules")
	rulesGroup.Get("/", ruleHandler.ListRules)
	rulesGroup.Get("/:id", ruleHandler.GetRule)

	// Rule modification requires rule_editor or admin role
	rulesProtected := rulesGroup.Group("", middleware.RequireAnyRole("admin", "rule_editor"))
	rulesProtected.Post("/", ruleHandler.CreateRule)
	rulesProtected.Put("/:id", ruleHandler.UpdateRule)
	rulesProtected.Delete("/:id", ruleHandler.DeleteRule)

	// Permissions management (admin only)
	permissionsGroup := v1.Group("/permissions", middleware.RequireRole("admin"))
	permissionsGroup.Get("/users", permissionsHandler.ListUsers)
	permissionsGroup.Get("/users/:id", permissionsHandler.GetUser)
	permissionsGroup.Put("/users/:id/active", permissionsHandler.UpdateUserActive)
	permissionsGroup.Post("/users/:userId/roles", permissionsHandler.AssignRole)
	permissionsGroup.Delete("/users/:userId/roles/:roleId", permissionsHandler.RemoveRole)
	permissionsGroup.Get("/roles", permissionsHandler.ListRoles)
	permissionsGroup.Get("/groups", permissionsHandler.ListGroups)
}
