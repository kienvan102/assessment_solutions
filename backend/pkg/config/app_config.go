package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds all centralized configuration variables for the application.
type AppConfig struct {
	// Server Port
	Port string

	// Payment System Configurations
	PaymentProcessingDelay time.Duration

	// Worker Pool Configurations
	WorkerPoolSize int

	// Code Review Security Configurations
	CodeReviewMaxBodySize int64
}

// LoadConfig reads from a .env file (if it exists) and parses environment variables.
// It provides sensible defaults if the variables are not set.
func LoadConfig() AppConfig {
	// Attempt to load .env file. It's okay if it doesn't exist (e.g. in production Docker)
	_ = godotenv.Load()

	cfg := AppConfig{
		Port:                   getEnv("PORT", "8080"),
		PaymentProcessingDelay: getEnvAsDuration("PAYMENT_PROCESSING_DELAY", 60*time.Second),
		WorkerPoolSize:         getEnvAsInt("WORKER_POOL_SIZE", 5),
		CodeReviewMaxBodySize:  getEnvAsInt64("CODE_REVIEW_MAX_BODY_SIZE", 1024*1024), // Default 1MB
	}

	return cfg
}

// getEnv fetches a string environment variable or returns a fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt fetches an integer environment variable or returns a fallback
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

// getEnvAsInt64 fetches a 64-bit integer environment variable or returns a fallback
func getEnvAsInt64(key string, fallback int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return fallback
}

// getEnvAsDuration fetches a duration string environment variable or returns a fallback
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return fallback
}
