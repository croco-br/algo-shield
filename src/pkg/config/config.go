package config

import (
	"fmt"
	"os"
	"strconv"

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
	Host string
	Port int
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

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "algoshield"),
			Password: getEnv("POSTGRES_PASSWORD", "algoshield_secret"),
			Database: getEnv("POSTGRES_DB", "algoshield"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvInt("REDIS_PORT", 6379),
		},
		API: APIConfig{
			Host: getEnv("API_HOST", "0.0.0.0"),
			Port: getEnvInt("API_PORT", 8080),
		},
		Worker: WorkerConfig{
			Concurrency: getEnvInt("WORKER_CONCURRENCY", 10),
			BatchSize:   getEnvInt("WORKER_BATCH_SIZE", 50),
		},
		General: GeneralConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
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
