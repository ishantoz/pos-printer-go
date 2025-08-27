package lib

import (
	"fmt"
	"time"

	"github.com/karalabe/hid"
)

type HIDWriter struct {
	device *hid.Device
}

// NewHIDWriter creates a new HID writer for thermal printer
func NewHIDWriter(vid, pid string) (*HIDWriter, error) {
	// Convert hex strings to integers
	var vendorID, productID uint16
	_, err := fmt.Sscanf(vid, "%x", &vendorID)
	if err != nil {
		return nil, fmt.Errorf("invalid vendor ID: %v", err)
	}
	_, err = fmt.Sscanf(pid, "%x", &productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Enumerate HID devices
	devices := hid.Enumerate(vendorID, productID)
	if len(devices) == 0 {
		return nil, fmt.Errorf("no HID device found with VID %s PID %s", vid, pid)
	}

	// Open the first matching device
	device, err := devices[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open HID device: %v", err)
	}

	return &HIDWriter{device: device}, nil
}

func (w *HIDWriter) Write(data []byte) (int, error) {
	// Thermal printers often expect data in specific packet sizes
	// Common packet sizes are 64 bytes (including report ID)
	packetSize := 64
	totalWritten := 0

	for i := 0; i < len(data); i += packetSize - 1 {
		end := i + packetSize - 1
		if end > len(data) {
			end = len(data)
		}

		// Create packet with report ID (usually 0x00 or 0x02 for printers)
		packet := make([]byte, packetSize)
		packet[0] = 0x00 // Report ID - try 0x00, 0x01, 0x02 if this doesn't work
		copy(packet[1:], data[i:end])

		n, err := w.device.Write(packet)
		if err != nil {
			return totalWritten, err
		}
		totalWritten += n - 1 // Subtract report ID byte

		// Small delay between packets (some printers need this)
		time.Sleep(10 * time.Millisecond)
	}

	return totalWritten, nil
}

func (w *HIDWriter) Close() error {
	if w.device != nil {
		return w.device.Close()
	}
	return nil
}
