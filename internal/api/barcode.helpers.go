package api

import "pos-printer/internal/model"

func (server *Server) applyDefaultsBarcodeHelper(req *model.PrintBarcodeRequest) {
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
	} else if req.PrintCount > server.cfg.PrinterConfig.MaxPrintCount {
		req.PrintCount = server.cfg.PrinterConfig.MaxPrintCount
	}
	if len(req.TopText) > server.cfg.PrinterConfig.MaxTopTextLength {
		req.TopText = req.TopText[:server.cfg.PrinterConfig.MaxTopTextLength]
	}
	if req.LabelGap.Length == 0 {
		req.LabelGap.Length = 0
	}
	if req.LabelGap.Offset == 0 {
		req.LabelGap.Offset = 0
	}
}
