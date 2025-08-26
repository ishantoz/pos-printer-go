package tspl

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/gousb"
)

func PrintBarcodeLabel(
	vidHexStr, pidHexStr string,
	sizeX, sizeY,
	dir int, topText,
	barcodeData string, printCount,
	gapLength, gapOffset int) error {

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

	ctx := gousb.NewContext()
	defer ctx.Close()

	dev, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		return fmt.Errorf("could not open device %04x:%04x: %w", vid, pid, err)
	}
	if dev == nil {
		return fmt.Errorf("printer %04x:%04x not found", vid, pid)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

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

	ep, err := intf.OutEndpoint(1)
	if err != nil {
		return fmt.Errorf("could not open endpoint: %w", err)
	}

	if gapLength == 0 {
		autodetectCmd := "AUTODETECT\r\n"
		if _, err := ep.Write([]byte(autodetectCmd)); err != nil {
			log.Printf("AUTODETECT failed to send: %v — falling back to default 2mm", err)
			gapLength = 2
			gapOffset = 0
		} else {
			time.Sleep(1500 * time.Millisecond)
		}
	}

	heightDots := sizeY * 8
	barcodeHeight := 70
	textHeight := 12
	spacing := 10
	totalBlock := textHeight + barcodeHeight + spacing
	yOffset := (heightDots - totalBlock) / 2

	tspl := ""
	tspl += fmt.Sprintf(
		"SIZE %d mm, %d mm\r\n",
		sizeX,
		sizeY,
	)
	if gapLength > 0 {
		tspl += fmt.Sprintf(
			"GAP %d mm, %d mm\r\n",
			gapLength,
			gapOffset,
		)
	}
	tspl += fmt.Sprintf("DIRECTION %d\r\n", dir)

	tspl += "CLS\r\n"

	tspl += "SET PRINTER DT\r\n"

	tspl += fmt.Sprintf(
		"TEXT 10,%d,\"1\",0,1,1,\"%s\"\r\n",
		yOffset,
		topText,
	)
	tspl += fmt.Sprintf(
		"BARCODE 0,%d,\"128\",%d,1,0,1,2,\"%s\"\r\n",
		yOffset+textHeight+spacing,
		barcodeHeight,
		barcodeData,
	)
	tspl += fmt.Sprintf(
		"PRINT %d,1\r\n",
		printCount,
	)
	tspl += "CUT\r\n"

	if _, err := ep.Write([]byte(tspl)); err != nil {
		return fmt.Errorf("failed to write TSPL data: %w", err)
	}

	return nil
}
