/*
*	pos-printer is a service that prints barcodes and receipts to a POS printer.
*		Features:
*			- Prints barcodes and receipts to a POS printer.
*			- Uses a SQLite database to store jobs.
*			- Uses a gousb library to communicate with the printer.
*			- Has an API that can be used to enqueue jobs and get the status of the jobs.
*			- Has a worker that processes jobs in the background.
*			- Has a server that can be used to enqueue jobs and get the status of the jobs.
*			- Has a server that can be used to print barcodes and receipts.
*			- Has a server that can be used to print receipts.
*			- Has a server that can be used to print barcodes and receipts.
*			- Has a server that can be used to print receipts.
*			- Has a server that can be used to print barcodes and receipts.
*		Implementation:
*			- Highly uses dependency injection to make the code more testable.
*			- Uses a gousb library to communicate with the printer.
*			- Uses a SQLite database to store jobs.
*			- Uses a gousb library to communicate with the printer.
*			- Uses a gousb library to communicate with the printer.
**/

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pos-printer/internal/api"
	"pos-printer/internal/config"
	"pos-printer/internal/db"
	"pos-printer/internal/job"
	"pos-printer/internal/printer"
)

func main() {
	cfg := config.Load()

	sqlite, err := db.NewSQLite(cfg)
	if err != nil {
		log.Fatalf("failed to initialize SQLite database: %v", err)
	}

	posPrinter := printer.NewPosPrinter()

	processor := job.NewProcessor(posPrinter, sqlite, cfg)
	processor.StartWorkers()

	server := api.NewServer(cfg, sqlite, posPrinter)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartTLS(); err != nil {
			log.Printf("server failed: %v", err)
		}
	}()

	fmt.Printf("🚀 POS Printer Service started securely on https://localhost%s\n", cfg.ServerConfig.Endpoint)

	<-sigChan
	fmt.Println("\n🛑 Shutting down gracefully...")

	posPrinter.Cleanup()
	fmt.Println("✅ Cleanup completed")
}
