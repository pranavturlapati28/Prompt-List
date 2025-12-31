package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DatabaseURL string
	Port        string
	Environment string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	// This is useful for local development
	_ = godotenv.Load()

	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://promptuser:promptpass@localhost:5432/prompttree?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}