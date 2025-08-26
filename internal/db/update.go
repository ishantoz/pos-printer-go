package db

import (
	"fmt"
	"log"
)

func (s *SQLite) UpdateStaleBarcodeJobs() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	_, err := s.db.Exec(
		`UPDATE barcode_jobs
			 SET status = ?, updatedAt = CURRENT_TIMESTAMP
			 WHERE status = ?
			 AND updatedAt < DATETIME('now', ?)`,
		s.cfg.WorkerConfig.JobStatus.StatusPending,
		s.cfg.WorkerConfig.JobStatus.StatusInProgress,
		fmt.Sprintf(
			"-%d minutes",
			int(s.cfg.WorkerConfig.StaleThreshold.Minutes()),
		),
	)
	return err
}

func (s *SQLite) UpdateBarcodeJobStatus(jobID int, status string) error {
	dbMu.Lock()
	defer dbMu.Unlock()
	_, err := s.db.Exec(
		`UPDATE barcode_jobs SET status = ?, updatedAt = CURRENT_TIMESTAMP WHERE id = ?`,
		status, jobID,
	)
	if err != nil {
		log.Printf("Error updating barcode job status: %v", err)
		return err
	}
	return nil
}

func (s *SQLite) UpdateBarcodeJobAttempts(jobID int) error {
	dbMu.Lock()
	defer dbMu.Unlock()
	_, err := s.db.Exec(
		`UPDATE barcode_jobs SET status = ?, attempts = attempts + 1, updatedAt = CURRENT_TIMESTAMP WHERE id = ?`,
		s.cfg.WorkerConfig.JobStatus.StatusInProgress, jobID,
	)
	if err != nil {
		return err
	}
	return nil
}
