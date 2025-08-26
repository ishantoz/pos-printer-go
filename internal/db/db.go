package db

import (
	"log"
	"pos-printer/internal/config"
	"pos-printer/internal/db/sqlite"

	"fmt"
)

type SQLiteContext struct {
	SQLitePath string
	Migrate    bool
}

func InitSQLite(dbContext *SQLiteContext) (*sqlite.SQLite, error) {
	db, err := sqlite.NewSQLite(dbContext.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	if dbContext.Migrate {
		if err := db.Migrate(); err != nil {
			return nil, fmt.Errorf("failed to migrate SQLite database: %v", err)
		}
	}

	return db, nil
}

func InitDB(cfg *config.DBConfig) {
	dbContext := &SQLiteContext{
		SQLitePath: cfg.SQLitePath,
		Migrate:    cfg.Migrate,
	}
	// Initialize SQLite database and migrate if needed
	_, err := InitSQLite(dbContext)
	if err != nil {
		log.Fatalf("failed to open SQLite database: %v", err)
	}
}
