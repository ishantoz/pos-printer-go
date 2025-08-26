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
	SQLitePath string
	Migrate    bool
}

type BarcodeConfig struct {
	MinSizeMM      int
	MaxSizeMM      int
	MinGapMM       int
	MaxGapMM       int
	MinGapOffsetMM int
	MaxGapOffsetMM int
	MinDirection   int
	MaxDirection   int
}

type PrinterConfig struct {
	MaxPrintCount        int
	MaxBarcodeDataLength int
	MaxTopTextLength     int
	BarcodeConfig        BarcodeConfig
}

type JobStatus struct {
	StatusPending    string
	StatusInProgress string
	StatusFailed     string
	StatusDone       string
}

type WorkerConfig struct {
	MaxJobAttempts     int
	BarcodeWorkerCount int
	JobStatus          JobStatus
	StaleThreshold     time.Duration
	StaleInterval      time.Duration
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
			Timeout:  10 * time.Second,
			CertPath: GetEnv("SERVER_CERT_PATH", "./certs/cert.pem"),
			KeyPath:  GetEnv("SERVER_KEY_PATH", "./certs/cert.key"),
		},
		DBConfig: DBConfig{
			SQLitePath: GetEnv("DB_SQLITE_PATH", "./data/db/pos-printer.sqlite.db"),
			Migrate:    GetEnvInt("DB_MIGRATE", 1) == 1,
		},
		PrinterConfig: PrinterConfig{
			MaxPrintCount:        GetEnvInt("MAX_BARCODE_PRINT_COUNT", 1000),
			MaxBarcodeDataLength: GetEnvInt("MAX_BARCODE_DATA_LENGTH", 100),
			MaxTopTextLength:     GetEnvInt("MAX_TOP_TEXT_LENGTH", 50),
			BarcodeConfig: BarcodeConfig{
				MinSizeMM:      5,
				MaxSizeMM:      200,
				MinGapMM:       0, // 0 means auto-detect
				MaxGapMM:       50,
				MinGapOffsetMM: -10,
				MaxGapOffsetMM: 10,
				MinDirection:   0,
				MaxDirection:   1,
			},
		},
		WorkerConfig: WorkerConfig{
			MaxJobAttempts:     GetEnvInt("MAX_JOB_ATTEMPTS", 3),
			BarcodeWorkerCount: GetEnvInt("BARCODE_WORKER_COUNT", 3),
			StaleThreshold:     time.Duration(10) * time.Minute,
			StaleInterval:      time.Duration(5) * time.Minute,
			JobStatus: JobStatus{
				StatusPending:    "pending",
				StatusInProgress: "in_progress",
				StatusFailed:     "failed",
				StatusDone:       "done",
			},
		},
	}
}
