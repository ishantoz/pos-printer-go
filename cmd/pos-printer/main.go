package main

import (
	"fmt"
	"log"

	"pos-printer/internal/api"
	"pos-printer/internal/config"
	"pos-printer/internal/db"
	// "pos-printer/internal/job"
	// "pos-printer/internal/printer"
)

func main() {

	cfg := config.Load()

	// Initialize SQLite database
	db.InitDB(&cfg.DBConfig)

	// Initialize TSPL printer
	// tsplPrinter := printer.NewTSPLPrinter(cfg.PrinterConfig)
	// processor := job.NewProcessor(db, tsplPrinter)

	// go processor.StartWorkers(cfg.WorkerConfig.WorkerCount)

	server := api.NewServer()

	if err := server.StartTLS(&cfg.ServerConfig); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	fmt.Printf("🚀 POS Printer Service started securely on https://localhost%s\n", cfg.ServerConfig.Endpoint)

}
