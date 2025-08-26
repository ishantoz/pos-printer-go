package db

import (
	"fmt"
	"pos-printer/internal/model"
	"time"
)

func (s *SQLite) EnqueueBarcodeJob(req model.PrintBarcodeRequest) (int64, error) {
	now := time.Now()
	dbMu.Lock()
	defer dbMu.Unlock()
	res, err := s.db.Exec(
		`INSERT INTO barcode_jobs 
		(vid, pid, sizeX, sizeY, direction, topText, barcodeData, printCount, labelGapLength, labelGapOffset, status, attempts, createdAt, updatedAt)
		 VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		req.VID, req.PID, req.SizeX, req.SizeY,
		req.Direction, req.TopText, req.BarcodeData,
		req.PrintCount, req.LabelGap.Length, req.LabelGap.Offset,
		"pending", 0, now, now,
	)
	if err != nil {
		fmt.Println("Failed to enqueue job", err)
		return 0, err
	}
	return res.LastInsertId()
}
