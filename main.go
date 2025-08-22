package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"pos-printer/lib/device"
	"pos-printer/lib/tspl"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MaxPrintCount        = 1000
	MaxBarcodeDataLength = 100
	MaxTopTextLength     = 50
	MaxJobAttempts       = 3
	WorkerCount          = 3
	DBPath               = "jobs.db"

	StaleThreshold = 10 * time.Minute
)

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusFailed     = "failed"
	StatusDone       = "done"
)

type PrintBarcodeRequest struct {
	VID         string `json:"vid"`
	PID         string `json:"pid"`
	SizeX       int    `json:"sizeX"`
	SizeY       int    `json:"sizeY"`
	Direction   int    `json:"direction"`
	TopText     string `json:"topText"`
	BarcodeData string `json:"barcodeData"`
	PrintCount  int    `json:"printCount"`
}

type Job struct {
	ID        int
	Request   PrintBarcodeRequest
	Status    string
	Attempts  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	db   *sql.DB
	dbMu sync.Mutex
)

func main() {
	if err := initDB(); err != nil {
		log.Fatalf("DB init error: %v", err)
	}

	go requeueStaleBarcodeJobs()

	for i := 0; i < WorkerCount; i++ {
		go workerBarcode(i + 1)
	}

	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	fmt.Println("🚀 POS Printer Service started securely on https://localhost:5000")

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Barcode Print API
	e.POST("/print-barcode-labels", enqueueBarcodeHandler)
	e.GET("/barcode-job-status/:id", jobBarcodeStatusHandler)

	// Start HTTPS server
	certPath := "./certs/cert.pem"
	keyPath := "./certs/cert.key"
	log.Printf("Starting HTTPS server on :5000")
	if err := e.StartTLS(":5000", certPath, keyPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("db ping error: %w", err)
	}
	barcodeJobTableStmt := `CREATE TABLE IF NOT EXISTS barcode_jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vid TEXT, pid TEXT,
		sizeX INTEGER, sizeY INTEGER,
		direction INTEGER, topText TEXT,
		barcodeData TEXT, printCount INTEGER,
		status TEXT, attempts INTEGER,
		createdAt DATETIME, updatedAt DATETIME
	);`

	receiptPDFJobTableStmt := `CREATE TABLE IF NOT EXISTS receipt_pdf_jobs (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path        TEXT    NOT NULL,
		connection_type  TEXT    NOT NULL CHECK(connection_type IN ('network','usb')),
		printer_ip       TEXT,
		printer_port     INTEGER,
		usb_vendor_id    INTEGER,
		usb_product_id   INTEGER,
		usb_interface    INTEGER DEFAULT 0,
		printer_width    INTEGER DEFAULT 576,
		threshold        INTEGER DEFAULT 100,
		feed_lines       INTEGER DEFAULT 1,
		zoom             REAL    DEFAULT 2.0,
		status           TEXT    DEFAULT 'pending',
		retry_count      INTEGER DEFAULT 0,
		last_error       TEXT,
		created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );`

	_, err = db.Exec(barcodeJobTableStmt)
	if err != nil {
		return err
	}
	_, err = db.Exec(receiptPDFJobTableStmt)
	if err != nil {
		return err
	}
	return err
}

func requeueStaleBarcodeJobs() {
	for {
		dbMu.Lock()
		_, err := db.Exec(
			`UPDATE barcode_jobs
			 SET status = ?, updatedAt = CURRENT_TIMESTAMP
			 WHERE status = ?
			   AND updatedAt < DATETIME('now', ?)`,
			StatusPending, StatusInProgress, fmt.Sprintf("-%d minutes", int(StaleThreshold.Minutes())),
		)
		dbMu.Unlock()
		if err != nil {
			log.Printf("Error requeuing stale jobs: %v", err)
		}
		time.Sleep(5 * time.Minute)
	}
}

func enqueueBarcodeHandler(c echo.Context) error {
	var req PrintBarcodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	applyBarcodeDefaults(&req)
	if err := validateBarcodeRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := device.CheckPrinter(req.VID, req.PID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Printer device not found, please check connected or not: %s", err)})
	}

	now := time.Now()
	dbMu.Lock()
	res, err := db.Exec(
		`INSERT INTO barcode_jobs (vid,pid,sizeX,sizeY,direction,topText,barcodeData,printCount,status,attempts,createdAt,updatedAt)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		req.VID, req.PID, req.SizeX, req.SizeY,
		req.Direction, req.TopText, req.BarcodeData,
		req.PrintCount, StatusPending, 0, now, now,
	)
	dbMu.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to enqueue job"})
	}
	id, _ := res.LastInsertId()
	return c.JSON(http.StatusAccepted, echo.Map{"jobId": id, "status": StatusPending})
}

func jobBarcodeStatusHandler(c echo.Context) error {
	id := c.Param("id")
	var status string
	err := db.QueryRow(`SELECT status FROM barcode_jobs WHERE id = ?`, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error fetching job status"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": status})
}

func applyBarcodeDefaults(req *PrintBarcodeRequest) {
	if req.VID == "" {
		req.VID = "0x0fe6"
	}
	if req.PID == "" {
		req.PID = "0x8800"
	}
	if req.SizeX == 0 {
		req.SizeX = 45
	}
	if req.SizeY == 0 {
		req.SizeY = 35
	}
	if req.PrintCount < 1 {
		req.PrintCount = 1
	} else if req.PrintCount > MaxPrintCount {
		req.PrintCount = MaxPrintCount
	}
	if len(req.TopText) > MaxTopTextLength {
		req.TopText = req.TopText[:MaxTopTextLength]
	}
}

func validateBarcodeRequest(req *PrintBarcodeRequest) error {
	if req.BarcodeData == "" {
		return errors.New("barcodeData is required")
	}
	if len(req.BarcodeData) > MaxBarcodeDataLength {
		return fmt.Errorf("barcodeData must not exceed %d chars", MaxBarcodeDataLength)
	}
	return nil
}

func fetchBarcodeJob() (*Job, error) {
	dbMu.Lock()
	defer dbMu.Unlock()
	row := db.QueryRow(`
		SELECT id, vid, pid, sizeX, sizeY, direction, topText, barcodeData, printCount, attempts
		FROM barcode_jobs WHERE status = ? AND attempts < ? ORDER BY createdAt LIMIT 1`,
		StatusPending, MaxJobAttempts,
	)

	var job Job
	var attempts int
	err := row.Scan(
		&job.ID,
		&job.Request.VID, &job.Request.PID,
		&job.Request.SizeX, &job.Request.SizeY, &job.Request.Direction,
		&job.Request.TopText, &job.Request.BarcodeData,
		&job.Request.PrintCount, &attempts,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	_, err = db.Exec(
		`UPDATE barcode_jobs SET status = ?, attempts = attempts + 1, updatedAt = CURRENT_TIMESTAMP WHERE id = ?`,
		StatusInProgress, job.ID,
	)
	if err != nil {
		return nil, err
	}

	job.Status = StatusInProgress
	job.Attempts = attempts + 1
	return &job, nil
}

func processBarcodeJob(workerID int, job *Job) {
	log.Printf("Worker %d processing job %d (attempt %d)", workerID, job.ID, job.Attempts)
	err := tspl.PrintBarcodeLabel(
		job.Request.VID, job.Request.PID,
		job.Request.SizeX, job.Request.SizeY,
		job.Request.Direction, job.Request.TopText,
		job.Request.BarcodeData, job.Request.PrintCount,
	)

	var newStatus string
	if err != nil {
		log.Printf("Worker %d job %d failed: %v", workerID, job.ID, err)
		if job.Attempts >= MaxJobAttempts {
			newStatus = StatusFailed
		} else {
			newStatus = StatusPending
		}
	} else {
		log.Printf("Worker %d job %d done", workerID, job.ID)
		newStatus = StatusDone
	}

	_, uerr := db.Exec(
		`UPDATE barcode_jobs SET status = ?, updatedAt = CURRENT_TIMESTAMP WHERE id = ?`,
		newStatus, job.ID,
	)
	if uerr != nil {
		log.Printf("Worker %d update job %d error: %v", workerID, job.ID, uerr)
	}
}

// Worker for printing jobs
func workerBarcode(id int) {
	for {
		job, err := fetchBarcodeJob()
		if err != nil {
			log.Printf("Worker %d: fetch error: %v", id, err)
			time.Sleep(time.Second)
			continue
		}
		if job == nil {
			time.Sleep(time.Second)
			continue
		}
		processBarcodeJob(id, job)
	}
}
