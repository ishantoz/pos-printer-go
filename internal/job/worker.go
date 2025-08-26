package job

import (
	"log"
	"time"
)

// Worker for printing jobs
func (p *Processor) workerBarcode(id int) {
	for {
		job, err := p.db.FetchBarcodeAndUpdateStatusToInProgress()
		if err != nil {
			log.Printf("Worker %d: fetch error: %v", id, err)
			time.Sleep(time.Second)
			continue
		}
		if job == nil {
			time.Sleep(time.Second)
			continue
		}
		p.processBarcodeJob(id, job)
	}
}

func (p *Processor) StartWorkers() {

	// Start requeue stale barcode jobs
	go p.RequeueStaleBarcodeJobs()

	// Start barcode workers
	for i := 0; i < p.cfg.WorkerConfig.BarcodeWorkerCount; i++ {
		go p.workerBarcode(i)
	}
}
