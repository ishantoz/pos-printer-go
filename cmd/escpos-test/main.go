package main

import (
	"log"
	"pos-printer/internal/printer"
)

// ESC/POS commands
var (
	ESC_INIT    = []byte{0x1B, 0x40}
	ESC_ALIGN_L = []byte{0x1B, 0x61, 0x00}
	CUT_FULL    = []byte{0x1D, 0x56, 0x00}
)

func ESC_FEED_N(n byte) []byte {
	return []byte{0x1B, 0x64, n}
}

func main() {
	vid := "0x0fe6"
	pid := "0x811e"

	printer := printer.NewPosPrinter()
	if err := printer.CheckPrinter(vid, pid); err != nil {
		log.Fatalf("Failed to check printer: %v", err)
	}

	writer, dev, err := printer.GetESCPOSWriter(vid, pid)
	if err != nil {
		log.Fatalf("Failed to get ESC writer: %v", err)
	}

	writer.Write(ESC_INIT)
	writer.Write(ESC_ALIGN_L)
	writer.Write([]byte("Hello, World!"))
	writer.Write(ESC_FEED_N(5))
	writer.Write(CUT_FULL)

	defer dev.Close()
}
