package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

const prefix = "POS_PRINTER_"

// LoadEnv loads environment variables from a .env file if it exists
func LoadEnv(filename string) {
	absPath, _ := filepath.Abs(filename)

	if err := godotenv.Load(absPath); err == nil {
		log.Println("Environment variables loaded from", absPath)
	}
}

// GetEnv returns the value of the environment variable or the fallback if not set.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(prefix + key); v != "" {
		return v
	}
	return fallback
}

// GetEnvInt parses an integer from environment variable or returns fallback
func GetEnvInt(key string, fallback int) int {
	if v := os.Getenv(prefix + key); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			return val
		}
	}
	return fallback
}
