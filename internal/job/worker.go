package job

import (
	"log"
	"time"
)

func (p *Processor) StopWorkers() {
	close(p.stopChan)
	p.wg.Wait()
}

func (p *Processor) workerBarcode(id int) {
	for {
		select {
		case <-p.stopChan:
			log.Printf("Worker %d stopping", id)
			return
		default:
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
}

func (p *Processor) StartWorkers() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.RequeueStaleBarcodeJobs()
	}()

	for i := 0; i < p.cfg.WorkerConfig.BarcodeWorkerCount; i++ {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			p.workerBarcode(id)
		}(i)
	}
}
