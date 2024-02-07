package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represents tha app configuration
type Config struct {
	DatabaseURL string
	AppPort     string
	AppSecret   string
}

// Load config form env variables o .env file
func LoadConfig() (*Config, error) {
	// Load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loafing .env file, using default env variables.")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "localhost:5440"),
		AppPort:     getEnv("APP_PORT", "8080"),
		AppSecret:   getEnv("APP_SECRET", "sup3rs3cr3t0"),
	}, nil
}

// returns the value of a env variable or its default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
