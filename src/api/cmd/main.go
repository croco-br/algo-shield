package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/algo-shield/algo-shield/src/api/internal/routes"
	"github.com/algo-shield/algo-shield/src/pkg/config"
	"github.com/algo-shield/algo-shield/src/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresPool(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redis, err := database.NewRedisClient(cfg.GetRedisAddr())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	// Create Fiber app with optimized settings
	app := fiber.New(fiber.Config{
		Prefork:               false,
		ServerHeader:          "AlgoShield",
		AppName:               "AlgoShield API",
		DisableStartupMessage: false,
		EnablePrintRoutes:     cfg.General.Environment == "development",
		ReadTimeout:           0,
		WriteTimeout:          0,
		IdleTimeout:           0,
	})

	// Setup routes
	routes.Setup(app, db.Pool, redis.Client, cfg)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)
	if cfg.API.TLSEnable {
		log.Printf("Starting API server with TLS on %s", addr)
		if err := app.ListenTLS(addr, cfg.API.TLSCert, cfg.API.TLSKey); err != nil {
			log.Fatalf("Failed to start server with TLS: %v", err)
		}
	} else {
		if cfg.General.Environment == "production" {
			log.Fatalf("TLS is required in production environment. Set TLS_ENABLE=true, TLS_CERT_PATH, and TLS_KEY_PATH")
		}
		log.Printf("Starting API server on %s (HTTP only - TLS not enabled)", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
