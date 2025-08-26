package job

import (
	"pos-printer/internal/config"
	"pos-printer/internal/db"
	"pos-printer/internal/printer"
)

type Processor struct {
	posPrinter *printer.PosPrinter
	db         *db.SQLite
	cfg        *config.Config
}

func NewProcessor(posPrinter *printer.PosPrinter, db *db.SQLite, cfg *config.Config) *Processor {
	return &Processor{
		posPrinter: posPrinter,
		db:         db,
		cfg:        cfg,
	}
}
