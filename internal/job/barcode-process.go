package job

import (
	"log"
	"pos-printer/internal/model"
	"time"
)

func (p *Processor) RequeueStaleBarcodeJobs() {
	for {
		err := p.db.UpdateStaleBarcodeJobs()
		if err != nil {
			log.Printf("Error requeuing stale jobs: %v", err)
		}
		time.Sleep(p.cfg.WorkerConfig.StaleInterval)
	}
}

func (p *Processor) processBarcodeJob(workerID int, job *model.BarcodeJob) {
	log.Printf("Barcode Worker %d processing job %d (attempt %d)", workerID, job.ID, job.Attempts)

	err := p.posPrinter.PrintBarcode(
		job.VID, job.PID,
		job.SizeX, job.SizeY,
		job.Direction, job.TopText,
		job.BarcodeData, job.PrintCount,
		job.LabelGapLength, job.LabelGapOffset,
	)

	var newStatus string
	if err != nil {
		log.Printf("Barcode Worker %d job %d failed: %v", workerID, job.ID, err)
		if job.Attempts >= p.cfg.WorkerConfig.MaxJobAttempts {
			newStatus = p.cfg.WorkerConfig.JobStatus.StatusFailed
		} else {
			newStatus = p.cfg.WorkerConfig.JobStatus.StatusPending
		}
	} else {
		log.Printf("Barcode Worker %d job %d done", workerID, job.ID)
		newStatus = p.cfg.WorkerConfig.JobStatus.StatusDone
	}

	uerr := p.db.UpdateBarcodeJobStatus(job.ID, newStatus)

	if uerr != nil {
		log.Printf("Barcode Worker %d update job %d error: %v", workerID, job.ID, uerr)
		return
	}

	log.Printf("Barcode Worker %d job %d done", workerID, job.ID)
}
