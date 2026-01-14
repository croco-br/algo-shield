package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	Host             string
	Port             int
	TLSEnable        bool
	TLSCert          string // Path to TLS certificate file
	TLSKey           string // Path to TLS private key file
	CORSAllowOrigins string // Comma-separated list of allowed CORS origins (use "*" for development only)
	BodyLimit        int    // Maximum request body size in bytes (SECURITY: prevent DoS attacks)
	Timeouts         APITimeouts
	Cache            APICacheConfig
}

type WorkerConfig struct {
	Concurrency int
	BatchSize   int
	Timeouts    WorkerTimeouts
	Retry       RetryConfig
	Queue       QueueConfig
	RulesReload RulesReloadConfig
	Asynq       AsynqConfig
}

// AsynqConfigAccessor provides access to Asynq configuration for the queue package
func (wc *WorkerConfig) AsynqConfig() AsynqConfig {
	return wc.Asynq
}

type WorkerTimeouts struct {
	TransactionProcessing time.Duration
	RuleEvaluation        time.Duration
}

type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

type QueueConfig struct {
	PopTimeout time.Duration
}

type RulesReloadConfig struct {
	Interval time.Duration
}

type AsynqConfig struct {
	DefaultTimeout     time.Duration // Default task timeout
	DefaultRetention   time.Duration // Default task retention period
	CriticalTimeout    time.Duration // Timeout for critical priority tasks
	LowPriorityTimeout time.Duration // Timeout for low priority tasks
	ShutdownTimeout    time.Duration // Worker shutdown timeout
}

// GetDefaultTimeout implements AsynqConfigProvider interface from queue package
func (ac AsynqConfig) GetDefaultTimeout() time.Duration {
	return ac.DefaultTimeout
}

// GetDefaultRetention implements AsynqConfigProvider interface from queue package
func (ac AsynqConfig) GetDefaultRetention() time.Duration {
	return ac.DefaultRetention
}

// GetCriticalTimeout implements AsynqConfigProvider interface from queue package
func (ac AsynqConfig) GetCriticalTimeout() time.Duration {
	return ac.CriticalTimeout
}

// GetLowPriorityTimeout implements AsynqConfigProvider interface from queue package
func (ac AsynqConfig) GetLowPriorityTimeout() time.Duration {
	return ac.LowPriorityTimeout
}

type APITimeouts struct {
	HandlerTimeout time.Duration // Default timeout for API handlers
	HealthCheck    time.Duration // Timeout for health check endpoints
}

type APICacheConfig struct {
	DashboardTTL time.Duration // Dashboard metrics cache TTL
	BrandingTTL  time.Duration // Branding cache TTL
	SystemTTL    time.Duration // System config cache TTL
	RulesTTL     time.Duration // Rules cache TTL
}

type GeneralConfig struct {
	Environment string
	LogLevel    string
}

type AuthConfig struct {
	JWTSecret                 string
	JWTExpirationHours        int
	JWTRefreshExpirationHours int
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
			Host:             getEnv("API_HOST", "0.0.0.0"),
			Port:             getEnvInt("API_PORT", 8080),
			TLSEnable:        getEnv("TLS_ENABLE", "") == "true",
			TLSCert:          getEnv("TLS_CERT_PATH", ""),
			TLSKey:           getEnv("TLS_KEY_PATH", ""),
			CORSAllowOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
			BodyLimit:        getEnvInt("API_BODY_LIMIT", 4*1024*1024), // SECURITY: Default 4MB to prevent DoS
			Timeouts: APITimeouts{
				HandlerTimeout: getEnvDuration("API_TIMEOUT_HANDLER", 500*time.Millisecond),
				HealthCheck:    getEnvDuration("API_TIMEOUT_HEALTH", 2*time.Second),
			},
			Cache: APICacheConfig{
				DashboardTTL: getEnvDuration("API_CACHE_DASHBOARD_TTL", 30*time.Second),
				BrandingTTL:  getEnvDuration("API_CACHE_BRANDING_TTL", 10*time.Minute),
				SystemTTL:    getEnvDuration("API_CACHE_SYSTEM_TTL", 5*time.Minute),
				RulesTTL:     getEnvDuration("API_CACHE_RULES_TTL", 5*time.Minute),
			},
		},
		Worker: WorkerConfig{
			Concurrency: getEnvInt("WORKER_CONCURRENCY", 10),
			BatchSize:   getEnvInt("WORKER_BATCH_SIZE", 50),
			Timeouts: WorkerTimeouts{
				TransactionProcessing: getEnvDuration("WORKER_TIMEOUT_TRANSACTION_PROCESSING", 300*time.Millisecond),
				RuleEvaluation:        getEnvDuration("WORKER_TIMEOUT_RULE_EVALUATION", 300*time.Millisecond),
			},
			Retry: RetryConfig{
				MaxAttempts:  getEnvInt("WORKER_RETRY_MAX_ATTEMPTS", 3),
				InitialDelay: getEnvDuration("WORKER_RETRY_INITIAL_DELAY", 100*time.Millisecond),
				MaxDelay:     getEnvDuration("WORKER_RETRY_MAX_DELAY", 5*time.Second),
				Multiplier:   getEnvFloat("WORKER_RETRY_MULTIPLIER", 2.0),
			},
			Queue: QueueConfig{
				PopTimeout: getEnvDuration("WORKER_QUEUE_POP_TIMEOUT", 1*time.Second),
			},
			RulesReload: RulesReloadConfig{
				Interval: getEnvDuration("WORKER_RULES_RELOAD_INTERVAL", 10*time.Second),
			},
			Asynq: AsynqConfig{
				DefaultTimeout:     getEnvDuration("WORKER_ASYNQ_DEFAULT_TIMEOUT", 5*time.Minute),
				DefaultRetention:   getEnvDuration("WORKER_ASYNQ_DEFAULT_RETENTION", 24*time.Hour),
				CriticalTimeout:    getEnvDuration("WORKER_ASYNQ_CRITICAL_TIMEOUT", 30*time.Second),
				LowPriorityTimeout: getEnvDuration("WORKER_ASYNQ_LOW_PRIORITY_TIMEOUT", 10*time.Minute),
				ShutdownTimeout:    getEnvDuration("WORKER_ASYNQ_SHUTDOWN_TIMEOUT", 30*time.Second),
			},
		},
		General: GeneralConfig{
			Environment: environment,
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			JWTSecret:                 jwtSecret,
			JWTExpirationHours:        getEnvInt("JWT_EXPIRATION_HOURS", 24),
			JWTRefreshExpirationHours: getEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 168), // Default: 7 days (168 hours)
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

	// Validate CORS configuration
	if isProduction {
		// In production, CORS must NOT use wildcard "*" - must specify allowed origins
		if config.API.CORSAllowOrigins == "*" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS must not use wildcard '*' in production. Specify allowed origins (comma-separated) for security")
		}
		// Ensure at least one origin is specified
		if strings.TrimSpace(config.API.CORSAllowOrigins) == "" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS must be specified in production (e.g., https://yourdomain.com)")
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

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		// Try parsing as integer seconds for convenience
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
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
