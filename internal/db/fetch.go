package db

import (
	"database/sql"
	"errors"
	"fmt"
	"pos-printer/internal/model"
)

func (s *SQLite) FetchBarcodeJob(id string) (*model.BarcodeJob, error) {
	var job model.BarcodeJob

	query := `SELECT 
	id, vid, pid, sizeX, sizeY, direction, topText, barcodeData, 
    printCount, labelGapLength, labelGapOffset, status, attempts, createdAt, updatedAt
    FROM barcode_jobs
    WHERE id = ?`

	row := s.db.QueryRow(query, id)

	err := row.Scan(
		&job.ID, &job.VID, &job.PID, &job.SizeX, &job.SizeY,
		&job.Direction, &job.TopText, &job.BarcodeData,
		&job.PrintCount, &job.LabelGapLength, &job.LabelGapOffset,
		&job.Status, &job.Attempts, &job.CreatedAt, &job.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error fetching barcode job", err)
		return nil, err
	}
	return &job, nil
}

func (s *SQLite) FetchBarcodeAndUpdateStatusToInProgress() (*model.BarcodeJob, error) {

	query := `
		SELECT id, vid, pid, sizeX, sizeY, direction, topText, barcodeData, printCount, labelGapLength, labelGapOffset, attempts
		FROM barcode_jobs WHERE status = ? AND attempts < ? ORDER BY createdAt LIMIT 1`

	row := s.db.QueryRow(query,
		s.cfg.WorkerConfig.JobStatus.StatusPending,
		s.cfg.WorkerConfig.MaxJobAttempts,
	)

	var job model.BarcodeJob

	err := row.Scan(
		&job.ID,
		&job.VID,
		&job.PID,
		&job.SizeX,
		&job.SizeY,
		&job.Direction,
		&job.TopText,
		&job.BarcodeData,
		&job.PrintCount,
		&job.LabelGapLength,
		&job.LabelGapOffset,
		&job.Attempts,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	err = s.UpdateBarcodeJobAttempts(job.ID)
	if err != nil {
		return nil, err
	}

	job.Status = s.cfg.WorkerConfig.JobStatus.StatusInProgress
	job.Attempts = job.Attempts + 1

	return &job, nil
}
