package tspl

import (
	"fmt"
	"strconv"

	"github.com/google/gousb"
)

// printBarcodeLabel opens the USB device, claims the endpoint, and sends a TSPL barcode label.
// vidHexStr, pidHexStr: USB Vendor and Product IDs as hex strings (e.g., "0x0fe6")
// sizeX, sizeY: label dimensions in mm
// dir: print direction (0-3)
// topText: human-readable text above the barcode
// barcodeData: the data to encode in the barcode
func PrintBarcodeLabel(vidHexStr, pidHexStr string, sizeX, sizeY, dir int, topText, barcodeData string, printCount int) error {
	// Parse hex strings to uint16
	vid64, err := strconv.ParseUint(vidHexStr, 0, 16)
	if err != nil {
		return fmt.Errorf("invalid Vendor ID %q: %w", vidHexStr, err)
	}
	pid64, err := strconv.ParseUint(pidHexStr, 0, 16)
	if err != nil {
		return fmt.Errorf("invalid Product ID %q: %w", pidHexStr, err)
	}
	vid := gousb.ID(uint16(vid64))
	pid := gousb.ID(uint16(pid64))

	// Create USB context
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open device
	dev, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		return fmt.Errorf("could not open device %04x:%04x: %w", vid, pid, err)
	}
	if dev == nil {
		return fmt.Errorf("printer %04x:%04x not found", vid, pid)
	}
	defer dev.Close()

	// Detach kernel driver if needed
	dev.SetAutoDetach(true)

	// Set configuration and claim interface
	cfg, err := dev.Config(1)
	if err != nil {
		return fmt.Errorf("could not set config: %w", err)
	}
	defer cfg.Close()

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		return fmt.Errorf("could not claim interface: %w", err)
	}
	defer intf.Close()

	// Open OUT endpoint
	ep, err := intf.OutEndpoint(1)
	if err != nil {
		return fmt.Errorf("could not open endpoint: %w", err)
	}

	// Calculate positioning in dots (203 dpi ~8 dots/mm)
	heightDots := sizeY * 8
	barcodeHeight := 80 // fixed height in dots
	textHeight := 12    // approx font 2 height
	spacing := 10       // dots between text and barcode
	totalBlock := textHeight + barcodeHeight + spacing
	yOffset := (heightDots - totalBlock) / 2

	// Build TSPL command string
	label := fmt.Sprintf(
		"SIZE %d mm, %d mm\r\n"+
			"GAP 2 mm, 0 mm\r\n"+
			"DIRECTION %d\r\n"+
			"CLS\r\n"+
			"SET PRINTER DT\r\n"+
			"TEXT 15,%d,\"2\",0,1,1,\"%s\"\r\n"+
			"BARCODE 0,%d,\"128\",%d,1,0,2,2,\"%s\"\r\n"+
			"PRINT %d,1\r\n"+
			"CUT\r\n",
		sizeX,
		sizeY,
		dir,
		yOffset,
		topText,
		yOffset+textHeight+spacing,
		barcodeHeight,
		barcodeData,
		printCount,
	)

	// Send label to printer
	if _, err := ep.Write([]byte(label)); err != nil {
		return fmt.Errorf("failed to write TSPL data: %w", err)
	}
	return nil
}

