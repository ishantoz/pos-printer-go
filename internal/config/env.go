package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

const prefix = "POS_PRINTER_"

func LoadEnv(filename string) {
	absPath, _ := filepath.Abs(filename)

	if err := godotenv.Load(absPath); err == nil {
		log.Println("Environment variables loaded from", absPath)
	}
}

func GetEnv(key, fallback string) string {
	if v := os.Getenv(prefix + key); v != "" {
		return v
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if v := os.Getenv(prefix + key); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			return val
		}
	}
	return fallback
}
