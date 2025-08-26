package printer

import (
	"fmt"
	"log"
	"time"
)

func (p *PosPrinter) PrintBarcode(
	vidHexStr, pidHexStr string,
	sizeX, sizeY, dir int,
	topText, barcodeData string,
	printCount, gapLength, gapOffset int) error {
	dev, err := p.OpenPosPrinter(vidHexStr, pidHexStr)
	if err != nil {
		return err
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
		"TEXT 15,%d,\"2\",0,1,1,\"%s\"\r\n",
		yOffset,
		topText,
	)
	tspl += fmt.Sprintf(
		"BARCODE 0,%d,\"128\",%d,1,0,2,2,\"%s\"\r\n",
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
