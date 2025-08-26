package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Println("Opening SQLite database at:", absPath)

	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Println("Directory does not exist, creating:", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Println("SQLite database does not exist, creating new one")
		file, err := os.Create(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	return &SQLite{db: db}, nil
}

func (s *SQLite) Migrate() error {
	stmts := []string{
		BarcodeJobTableStmt,
		ReceiptPDFJobTableStmt,
	}

	executeStmt := func(stmt string) error {
		_, err := s.db.Exec(stmt)
		return err
	}

	for _, stmt := range stmts {
		if err := executeStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}
