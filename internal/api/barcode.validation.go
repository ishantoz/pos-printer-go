package api

import (
	"errors"
	"fmt"
	"pos-printer/internal/model"
	"strings"
)

func (server *Server) validateBarcodeRequest(req *model.PrintBarcodeRequest) error {

	printerConfig := server.cfg.PrinterConfig
	barcodeConfig := printerConfig.BarcodeConfig

	// required
	if strings.TrimSpace(req.BarcodeData) == "" {
		return errors.New("barcodeData is required")
	}
	if len(req.BarcodeData) > printerConfig.MaxBarcodeDataLength {
		return fmt.Errorf(
			"barcodeData must not exceed %d chars",
			printerConfig.MaxBarcodeDataLength,
		)
	}

	// sizes
	if req.SizeX < barcodeConfig.MinSizeMM || req.SizeX > barcodeConfig.MaxSizeMM {
		return fmt.Errorf(
			"sizeX must be between %d and %d mm",
			barcodeConfig.MinSizeMM,
			barcodeConfig.MaxSizeMM,
		)
	}
	if req.SizeY < barcodeConfig.MinSizeMM || req.SizeY > barcodeConfig.MaxSizeMM {
		return fmt.Errorf(
			"sizeY must be between %d and %d mm",
			barcodeConfig.MinSizeMM,
			barcodeConfig.MaxSizeMM,
		)
	}

	// direction
	if req.Direction < barcodeConfig.MinDirection || req.Direction > barcodeConfig.MaxDirection {
		return fmt.Errorf(
			"direction must be %d or %d",
			barcodeConfig.MinDirection,
			barcodeConfig.MaxDirection,
		)
	}

	// print count
	if req.PrintCount < 1 || req.PrintCount > printerConfig.MaxPrintCount {
		return fmt.Errorf(
			"printCount must be between 1 and %d",
			printerConfig.MaxPrintCount,
		)
	}

	// top text length
	if len(req.TopText) > printerConfig.MaxTopTextLength {
		return fmt.Errorf(
			"topText must not exceed %d characters",
			printerConfig.MaxTopTextLength,
		)
	}

	// labelGap validations (0 allowed => auto-detect)
	if req.LabelGap.Length < barcodeConfig.MinGapMM || req.LabelGap.Length > barcodeConfig.MaxGapMM {
		return fmt.Errorf(
			"labelGap.length must be between %d and %d mm (0 means auto-detect)",
			barcodeConfig.MinGapMM,
			barcodeConfig.MaxGapMM,
		)
	}
	if req.LabelGap.Offset < barcodeConfig.MinGapOffsetMM || req.LabelGap.Offset > barcodeConfig.MaxGapOffsetMM {
		return fmt.Errorf(
			"labelGap.offset must be between %d and %d mm",
			barcodeConfig.MinGapOffsetMM, barcodeConfig.MaxGapOffsetMM)
	}

	return nil
}
