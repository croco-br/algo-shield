package routes

import (
	"strings"

	"github.com/algo-shield/algo-shield/src/api/internal/auth"
	"github.com/algo-shield/algo-shield/src/api/internal/branding"
	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/health"
	"github.com/algo-shield/algo-shield/src/api/internal/permissions"
	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	"github.com/algo-shield/algo-shield/src/api/internal/rules"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/middleware"
	"github.com/algo-shield/algo-shield/src/api/internal/transactions"
	"github.com/algo-shield/algo-shield/src/api/internal/user"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	rulespkg "github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config) {
	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.SecurityHeaders()) // Security headers for Brave compatibility
	app.Use(middleware.CORS())

	// Create repositories (infrastructure layer - can create concrete types)
	roleRepo := roles.NewPostgresRepository(db)
	groupRepo := groups.NewPostgresRepository(db)
	userRepo := user.NewPostgresUserRepository(db)
	userTxManager := user.NewPostgresTransactionManager(db)
	permissionsUserRepo := permissions.NewPostgresUserRepository(db)
	transactionRepo := transactions.NewPostgresRepository(db)
	ruleRepo := rulespkg.NewPostgresRepository(db, redis)
	brandingRepo := branding.NewPostgresRepository(db, redis)

	// Create services with dependency injection (business layer - receives interfaces)
	roleService := roles.NewService(roleRepo)
	groupService := groups.NewService(groupRepo)
	userService := user.NewService(userRepo, roleRepo, userTxManager, roleService, groupService)
	authService := auth.NewService(cfg, userService)
	permissionsService := permissions.NewService(permissionsUserRepo, roleService, groupService)
	transactionService := transactions.NewService(transactionRepo, redis)
	brandingService := branding.NewService(brandingRepo)

	// Create handlers with dependency injection (presentation layer - receives interfaces)
	authHandler := auth.NewHandler(authService, userService)
	permissionsHandler := permissions.NewHandler(permissionsService)
	roleHandler := roles.NewHandler(roleService)
	groupHandler := groups.NewHandler(groupService)
	transactionHandler := transactions.NewHandler(transactionService)
	ruleHandler := rules.NewHandler(ruleRepo)
	healthHandler := health.NewHandler(db, redis)
	brandingHandler := branding.NewHandler(brandingService)

	// Health routes (public)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// Register public API endpoints before creating protected groups
	// These must be registered as specific routes to avoid being caught by the v1 middleware
	app.Post("/api/v1/auth/register", authHandler.Register)
	app.Post("/api/v1/auth/login", authHandler.Login)
	app.Get("/api/v1/branding", brandingHandler.GetBranding)

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
	permissionsGroup.Post("/users/:userId/roles", roleHandler.AssignRole)
	permissionsGroup.Delete("/users/:userId/roles/:roleId", roleHandler.RemoveRole)

	// Roles management (admin only)
	rolesGroup := v1.Group("/roles", middleware.RequireRole("admin"))
	rolesGroup.Get("/", roleHandler.ListRoles)
	rolesGroup.Get("/:id", roleHandler.GetRole)

	// Groups management (admin only)
	groupsGroup := v1.Group("/groups", middleware.RequireRole("admin"))
	groupsGroup.Get("/", groupHandler.ListGroups)
	groupsGroup.Get("/:id", groupHandler.GetGroup)

	// Branding management (admin only)
	v1.Put("/branding", middleware.RequireRole("admin"), brandingHandler.UpdateBranding)

	// 404 handler - always return JSON for API routes
	app.Use(func(c *fiber.Ctx) error {
		// Only handle 404 for API routes
		if strings.HasPrefix(c.Path(), "/api") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "API endpoint not found",
				"path":  c.Path(),
			})
		}
		// For non-API routes, return default 404 (useful for SPA routing)
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})
}
