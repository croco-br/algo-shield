package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("POSTGRES_HOST", "testhost")
	os.Setenv("POSTGRES_PORT", "5433")
	os.Setenv("API_PORT", "9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("Expected database host 'testhost', got '%s'", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Expected database port 5433, got %d", cfg.Database.Port)
	}

	if cfg.API.Port != 9090 {
		t.Errorf("Expected API port 9090, got %d", cfg.API.Port)
	}

	// Clean up
	os.Unsetenv("POSTGRES_HOST")
	os.Unsetenv("POSTGRES_PORT")
	os.Unsetenv("API_PORT")
}

func TestGetDatabaseDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			Database: "testdb",
		},
	}

	dsn := cfg.GetDatabaseDSN()
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

	if dsn != expected {
		t.Errorf("Expected DSN '%s', got '%s'", expected, dsn)
	}
}

func TestGetRedisAddr(t *testing.T) {
	cfg := &Config{
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
	}

	addr := cfg.GetRedisAddr()
	expected := "localhost:6379"

	if addr != expected {
		t.Errorf("Expected Redis address '%s', got '%s'", expected, addr)
	}
}

