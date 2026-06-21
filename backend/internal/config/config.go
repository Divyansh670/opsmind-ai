package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	GitHubWebhookSecret string
	GitHubToken         string
	GroqAPIKey          string
	GroqModelID         string
	MaxWorkers          int
}

func Load() (*Config, error) {
	maxWorkers, err := strconv.Atoi(getEnv("MAX_WORKERS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_WORKERS value: %w", err)
	}

	cfg := &Config{
		ServerPort:          getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5433"),
		DBUser:              getEnv("DB_USER", "opsmind_user"),
		DBPassword:          getEnv("DB_PASSWORD", "opsmind_pass_dev"),
		DBName:              getEnv("DB_NAME", "opsmind_db"),
		DBSSLMode:           getEnv("DB_SSL_MODE", "disable"),
		GitHubWebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
		GitHubToken:         getEnv("GITHUB_TOKEN", ""),
		GroqAPIKey:          getEnv("GROQ_API_KEY", ""),
		GroqModelID:         getEnv("GROQ_MODEL_ID", "llama-3.3-70b-versatile"),
		MaxWorkers:          maxWorkers,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

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

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
