package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all environment-driven configuration for the backend
type Config struct {
	// Server
	ServerPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// GitHub
	GitHubWebhookSecret string
	GitHubToken         string

	// LLM (we will use Groq free tier)
	GroqAPIKey  string
	GroqModelID string

	// Worker Pool
	MaxWorkers int
}

// Load reads all config from environment variables with sensible defaults
func Load() (*Config, error) {
	maxWorkers, err := strconv.Atoi(getEnv("MAX_WORKERS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_WORKERS value: %w", err)
	}

	cfg := &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "opsmind_user"),
		DBPassword: getEnv("DB_PASSWORD", "opsmind_pass_dev"),
		DBName:     getEnv("DB_NAME", "opsmind_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// GitHub
		GitHubWebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
		GitHubToken:         getEnv("GITHUB_TOKEN", ""),

		// LLM
		GroqAPIKey:  getEnv("GROQ_API_KEY", ""),
		GroqModelID: getEnv("GROQ_MODEL_ID", "llama3-70b-8192"),

		// Worker Pool
		MaxWorkers: maxWorkers,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that critical config values are present
func (c *Config) validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

// DSN returns the PostgreSQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

// getEnv reads an env variable, returning fallback if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
