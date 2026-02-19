package routes

import (
	"strings"

	"github.com/algo-shield/algo-shield/src/api/internal"
	"github.com/algo-shield/algo-shield/src/api/internal/auth"
	"github.com/algo-shield/algo-shield/src/api/internal/branding"
	"github.com/algo-shield/algo-shield/src/api/internal/dashboard"
	"github.com/algo-shield/algo-shield/src/api/internal/groups"
	"github.com/algo-shield/algo-shield/src/api/internal/health"
	"github.com/algo-shield/algo-shield/src/api/internal/permissions"
	"github.com/algo-shield/algo-shield/src/api/internal/roles"
	"github.com/algo-shield/algo-shield/src/api/internal/rules"
	"github.com/algo-shield/algo-shield/src/api/internal/schemas"
	"github.com/algo-shield/algo-shield/src/api/internal/shared/middleware"
	"github.com/algo-shield/algo-shield/src/api/internal/system"
	"github.com/algo-shield/algo-shield/src/api/internal/transactions"
	"github.com/algo-shield/algo-shield/src/api/internal/user"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/queue"
	rulespkg "github.com/algo-shield/algo-shield/src/pkg/rules"
	"github.com/algo-shield/algo-shield/src/pkg/tokenrevoke"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, asynqClient *queue.AsynqClient, cfg *config.Config) {
	// Set global config for timeout access
	internal.SetGlobalConfig(cfg)

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.SecurityHeaders()) // Security headers for Brave compatibility
	app.Use(middleware.CORS(cfg.API.CORSAllowOrigins))
	app.Use(middleware.SyntheticModeMiddleware()) // Extract synthetic mode from header

	// Create repositories (infrastructure layer - can create concrete types)
	roleRepo := roles.NewPostgresRepository(db)
	groupRepo := groups.NewPostgresRepository(db)
	userRepo := user.NewPostgresUserRepository(db)
	userTxManager := user.NewPostgresTransactionManager(db)
	permissionsUserRepo := permissions.NewPostgresUserRepository(db)
	transactionRepo := transactions.NewPostgresRepository(db)
	ruleRepo := rulespkg.NewPostgresRepository(db, redis, internal.GetCacheTTL("rules"))
	brandingRepo := branding.NewPostgresRepository(db, redis)
	schemaRepo := schemas.NewPostgresRepository(db, redis)
	dashboardRepo := dashboard.NewPostgresRepository(db)
	systemRepo := system.NewPostgresRepository(db, redis)

	// Create services with dependency injection (business layer - receives interfaces)
	roleService := roles.NewService(roleRepo)
	groupService := groups.NewService(groupRepo)
	userService := user.NewService(userRepo, roleRepo, userTxManager, roleService, groupService)
	tokenRevokeService := tokenrevoke.NewService(redis)
	authService := auth.NewService(cfg, userService, tokenRevokeService)
	permissionsService := permissions.NewService(permissionsUserRepo, roleService, groupService, authService)
	transactionService := transactions.NewService(transactionRepo, asynqClient)
	brandingService := branding.NewService(brandingRepo)
	// Wrap asynqClient in adapter for schema service
	schemaTaskEnqueuer := schemas.NewAsynqAdapter(asynqClient)
	schemaService := schemas.NewService(schemaRepo, schemaTaskEnqueuer)
	dashboardService := dashboard.NewService(dashboardRepo, redis)
	systemService := system.NewService(systemRepo)

	// Create handlers with dependency injection (presentation layer - receives interfaces)
	authHandler := auth.NewHandler(authService, userService, redis)
	permissionsHandler := permissions.NewHandler(permissionsService)
	roleHandler := roles.NewHandler(roleService)
	groupHandler := groups.NewHandler(groupService)
	transactionHandler := transactions.NewHandler(transactionService)
	ruleHandler := rules.NewHandler(ruleRepo)
	healthHandler := health.NewHandler(db, redis)
	brandingHandler := branding.NewHandler(brandingService)
	schemaHandler := schemas.NewHandler(schemaService)
	dashboardHandler := dashboard.NewHandler(dashboardService)
	systemHandler := system.NewHandler(systemService)

	// Health routes (public)
	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	// Register public API endpoints before creating protected groups
	// These must be registered as specific routes to avoid being caught by the v1 middleware
	// SECURITY: Apply rate limiting to auth endpoints to prevent brute force attacks
	app.Post("/api/v1/auth/register",
		middleware.RateLimiter(redis, middleware.RegisterRateLimit),
		authHandler.Register,
	)
	app.Post("/api/v1/auth/login",
		middleware.RateLimiter(redis, middleware.LoginRateLimit),
		authHandler.Login,
	)
	app.Post("/api/v1/auth/refresh",
		middleware.RateLimiter(redis, middleware.LoginRateLimit),
		authHandler.RefreshToken,
	)
	app.Get("/api/v1/branding", brandingHandler.GetBranding)

	// API v1 (protected)
	v1 := app.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(authHandler))

	// SECURITY: CSRF Protection for all state-changing requests (POST, PUT, PATCH, DELETE)
	// Excluded paths: login, register, refresh, branding (public GET endpoint)
	v1.Use(middleware.CSRFProtection(middleware.CSRFConfig{
		Redis: redis,
		ExcludedPaths: []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/refresh",
			"/api/v1/branding", // Public GET endpoint
		},
	}))

	// Current user
	v1.Get("/auth/me", authHandler.GetCurrentUser)
	v1.Post("/auth/logout", authHandler.Logout)

	// Transaction routes (protected)
	transactionsGroup := v1.Group("/transactions")
	transactionsGroup.Post("/", transactionHandler.ProcessTransaction)
	transactionsGroup.Get("/", transactionHandler.ListTransactions)
	transactionsGroup.Get("/:id", transactionHandler.GetTransaction)
	transactionsGroup.Patch("/:id/approve", transactionHandler.ApproveTransaction)
	transactionsGroup.Patch("/:id/reject", transactionHandler.RejectTransaction)

	// Dashboard routes (protected)
	v1.Get("/dashboard/metrics", dashboardHandler.GetMetrics)

	// Rule routes (protected)
	rulesGroup := v1.Group("/rules")
	rulesGroup.Get("/", ruleHandler.ListRules)
	rulesGroup.Get("/:id", ruleHandler.GetRule)

	// Rule modification requires rule_editor or admin role
	rulesProtected := rulesGroup.Group("", middleware.RequireAnyRole("admin", "rule_editor"))
	rulesProtected.Post("/", ruleHandler.CreateRule)
	rulesProtected.Put("/:id", ruleHandler.UpdateRule)
	rulesProtected.Delete("/:id", ruleHandler.DeleteRule)

	// Schema routes (protected)
	schemasGroup := v1.Group("/schemas")
	schemasGroup.Get("/", schemaHandler.ListSchemas)
	schemasGroup.Get("/:id", schemaHandler.GetSchema)

	// Schema modification requires rule_editor or admin role
	schemasProtected := schemasGroup.Group("", middleware.RequireAnyRole("admin", "rule_editor"))
	schemasProtected.Post("/", schemaHandler.CreateSchema)
	schemasProtected.Put("/:id", schemaHandler.UpdateSchema)
	schemasProtected.Delete("/:id", schemaHandler.DeleteSchema)
	schemasProtected.Post("/:id/parse", schemaHandler.ParseSchema)
	schemasProtected.Post("/:id/generate-events", schemaHandler.GenerateEvents)

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

	// System configuration routes (protected)
	systemGroup := v1.Group("/system")
	systemGroup.Get("/mode", systemHandler.GetSyntheticMode)
	systemGroup.Put("/mode", middleware.RequireRole("admin"), systemHandler.SetSyntheticMode)

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
