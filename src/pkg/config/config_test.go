package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables (including required secrets)
	_ = os.Setenv("POSTGRES_HOST", "testhost")
	_ = os.Setenv("POSTGRES_PORT", "5433")
	_ = os.Setenv("API_PORT", "9090")
	// Set required secrets with strong values for testing
	_ = os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-characters-long-for-validation")
	_ = os.Setenv("POSTGRES_PASSWORD", "test-db-password-minimum-16-chars")

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
	_ = os.Unsetenv("POSTGRES_HOST")
	_ = os.Unsetenv("POSTGRES_PORT")
	_ = os.Unsetenv("API_PORT")
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("POSTGRES_PASSWORD")
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

func TestLoad_ProductionRequiresTLS(t *testing.T) {
	// Set production environment
	_ = os.Setenv("ENVIRONMENT", "production")
	// Use strong secrets that pass production validation
	_ = os.Setenv("JWT_SECRET", "MyStr0ng!JWT$ecret#Key@2024#Minimum32Chars")
	_ = os.Setenv("POSTGRES_PASSWORD", "MyStr0ng!DB$Pass#2024")

	// Try to load config without TLS - should fail
	_, err := Load()
	if err == nil {
		t.Error("Expected error when loading production config without TLS, but got none")
		// Clean up if test failed
		_ = os.Unsetenv("ENVIRONMENT")
		_ = os.Unsetenv("JWT_SECRET")
		_ = os.Unsetenv("POSTGRES_PASSWORD")
		return
	}

	// Now set TLS properly
	_ = os.Setenv("TLS_ENABLE", "true")
	_ = os.Setenv("TLS_CERT_PATH", "/path/to/cert.pem")
	_ = os.Setenv("TLS_KEY_PATH", "/path/to/key.pem")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load production config with TLS: %v", err)
	}

	if !cfg.API.TLSEnable {
		t.Error("Expected TLS to be enabled in production")
	}

	if cfg.API.TLSCert == "" {
		t.Error("Expected TLS certificate path to be set")
	}

	if cfg.API.TLSKey == "" {
		t.Error("Expected TLS key path to be set")
	}

	// Clean up
	_ = os.Unsetenv("ENVIRONMENT")
	_ = os.Unsetenv("JWT_SECRET")
	_ = os.Unsetenv("POSTGRES_PASSWORD")
	_ = os.Unsetenv("TLS_ENABLE")
	_ = os.Unsetenv("TLS_CERT_PATH")
	_ = os.Unsetenv("TLS_KEY_PATH")
}
