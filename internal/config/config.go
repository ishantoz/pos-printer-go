package config

import (
	"time"
)

type ServerConfig struct {
	Endpoint string
	Timeout  time.Duration
	CertPath string
	KeyPath  string
}

type DBConfig struct {
	SQLitePath     string
	Migrate        bool
	StaleThreshold time.Duration
}

type PrinterConfig struct {
	MaxPrintCount        int
	MaxBarcodeDataLength int
	MaxTopTextLength     int
}

type WorkerConfig struct {
	MaxJobAttempts int
	WorkerCount    int
}

type Config struct {
	ServerConfig  ServerConfig
	DBConfig      DBConfig
	PrinterConfig PrinterConfig
	WorkerConfig  WorkerConfig
}

func Load() *Config {

	// Load .env file if it exists in the current directory process
	LoadEnv(".env")

	return &Config{
		ServerConfig: ServerConfig{
			Endpoint: GetEnv("ENDPOINT", ":5000"),
			Timeout:  time.Duration(GetEnvInt("SERVER_TIMEOUT", 10)) * time.Second,
			CertPath: GetEnv("SERVER_CERT_PATH", "./certs/cert.pem"),
			KeyPath:  GetEnv("SERVER_KEY_PATH", "./certs/cert.key"),
		},
		DBConfig: DBConfig{
			SQLitePath:     GetEnv("DB_SQLITE_PATH", "./data/db/pos-printer.sqlite.db"),
			Migrate:        GetEnvInt("DB_MIGRATE", 1) == 1,
			StaleThreshold: time.Duration(GetEnvInt("DB_STALE_THRESHOLD", 10)) * time.Minute,
		},
		PrinterConfig: PrinterConfig{
			MaxPrintCount:        GetEnvInt("MAX_PRINT_COUNT", 1000),
			MaxBarcodeDataLength: GetEnvInt("MAX_BARCODE_DATA_LENGTH", 100),
			MaxTopTextLength:     GetEnvInt("MAX_TOP_TEXT_LENGTH", 50),
		},
		WorkerConfig: WorkerConfig{
			MaxJobAttempts: GetEnvInt("MAX_JOB_ATTEMPTS", 3),
			WorkerCount:    GetEnvInt("MAX_WORKER_COUNT", 3),
		},
	}
}
