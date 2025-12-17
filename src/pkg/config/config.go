package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	API      APIConfig
	Worker   WorkerConfig
	General  GeneralConfig
	Auth     AuthConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Host string
	Port int
}

type APIConfig struct {
	Host      string
	Port      int
	TLSEnable bool
	TLSCert   string // Path to TLS certificate file
	TLSKey    string // Path to TLS private key file
}

type WorkerConfig struct {
	Concurrency int
	BatchSize   int
}

type GeneralConfig struct {
	Environment string
	LogLevel    string
}

type AuthConfig struct {
	JWTSecret          string
	JWTExpirationHours int
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	environment := getEnv("ENVIRONMENT", "development")
	isProduction := environment == "production"

	// Get JWT secret - REQUIRED, no defaults allowed
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required and must be set")
	}

	// Validate JWT secret strength
	if err := validateSecretStrength("JWT_SECRET", jwtSecret, isProduction, 32); err != nil {
		return nil, err
	}

	// Get database password - REQUIRED, no defaults allowed
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable is required and must be set")
	}

	// Validate database password strength
	if err := validateSecretStrength("POSTGRES_PASSWORD", dbPassword, isProduction, 16); err != nil {
		return nil, err
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "algoshield"),
			Password: dbPassword,
			Database: getEnv("POSTGRES_DB", "algoshield"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvInt("REDIS_PORT", 6379),
		},
		API: APIConfig{
			Host:      getEnv("API_HOST", "0.0.0.0"),
			Port:      getEnvInt("API_PORT", 8080),
			TLSEnable: getEnv("TLS_ENABLE", "") == "true",
			TLSCert:   getEnv("TLS_CERT_PATH", ""),
			TLSKey:    getEnv("TLS_KEY_PATH", ""),
		},
		Worker: WorkerConfig{
			Concurrency: getEnvInt("WORKER_CONCURRENCY", 10),
			BatchSize:   getEnvInt("WORKER_BATCH_SIZE", 50),
		},
		General: GeneralConfig{
			Environment: environment,
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			JWTSecret:          jwtSecret,
			JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
	}

	// Validate TLS configuration
	if isProduction {
		// In production, TLS is REQUIRED
		if !config.API.TLSEnable {
			return nil, fmt.Errorf("TLS_ENABLE=true is required in production environment for security")
		}
		if config.API.TLSCert == "" {
			return nil, fmt.Errorf("TLS_CERT_PATH is required when TLS_ENABLE=true (required in production)")
		}
		if config.API.TLSKey == "" {
			return nil, fmt.Errorf("TLS_KEY_PATH is required when TLS_ENABLE=true (required in production)")
		}
	} else if config.API.TLSEnable {
		// In development/test, if TLS is enabled, both cert and key must be provided
		if config.API.TLSCert == "" || config.API.TLSKey == "" {
			return nil, fmt.Errorf("both TLS_CERT_PATH and TLS_KEY_PATH must be provided when TLS_ENABLE=true")
		}
	}

	return config, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
}

func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// validateSecretStrength validates the strength of a secret
// - minLength: minimum required length
// - isProduction: if true, enforces stricter rules
func validateSecretStrength(secretName, secret string, isProduction bool, minLength int) error {
	// Check minimum length
	if len(secret) < minLength {
		return fmt.Errorf("%s must be at least %d characters long", secretName, minLength)
	}

	// Maximum length to prevent DoS attacks
	if len(secret) > 512 {
		return fmt.Errorf("%s must be at most 512 characters long", secretName)
	}

	// Production-specific validations
	if isProduction {
		// Check for common weak/default values
		weakSecrets := []string{
			"change-me-in-production",
			"change-me",
			"secret",
			"password",
			"algoshield_secret",
			"default",
			"test",
			"12345678",
		}
		secretLower := strings.ToLower(secret)
		for _, weak := range weakSecrets {
			if secretLower == weak || strings.Contains(secretLower, weak) {
				return fmt.Errorf("%s contains a weak or default value which is not allowed in production", secretName)
			}
		}

		// Enforce minimum complexity in production
		hasUpper := false
		hasLower := false
		hasDigit := false
		hasSpecial := false

		for _, char := range secret {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}

		// Require at least 3 out of 4 character types for strong secrets
		complexityCount := 0
		if hasUpper {
			complexityCount++
		}
		if hasLower {
			complexityCount++
		}
		if hasDigit {
			complexityCount++
		}
		if hasSpecial {
			complexityCount++
		}

		if complexityCount < 3 {
			return fmt.Errorf("%s must contain at least 3 of the following: uppercase letters, lowercase letters, digits, special characters (required in production)", secretName)
		}

		// Check for repeated characters (weak pattern)
		if hasRepeatedPattern(secret) {
			return fmt.Errorf("%s contains repeated patterns which weakens security (not allowed in production)", secretName)
		}
	}

	return nil
}

// hasRepeatedPattern checks if a string has obvious repeated patterns
func hasRepeatedPattern(s string) bool {
	if len(s) < 4 {
		return false
	}

	// Check for sequences like "aaaa", "1234", "abcd"
	for i := 0; i < len(s)-3; i++ {
		substr := s[i : i+4]
		// Check if all characters are the same
		allSame := true
		for j := 1; j < len(substr); j++ {
			if substr[j] != substr[0] {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}

		// Check for sequential patterns (simple check)
		isSequential := true
		for j := 1; j < len(substr); j++ {
			if substr[j] != substr[j-1]+1 && substr[j] != substr[j-1]-1 {
				isSequential = false
				break
			}
		}
		if isSequential {
			return true
		}
	}

	return false
}
