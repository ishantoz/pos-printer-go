package device

import (
	"fmt"
	"strconv"

	"github.com/google/gousb"
)

// CheckPrinter checks if the printer exists.
func CheckPrinter(vidHexStr, pidHexStr string) error {
	// Parse hex strings
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
		return fmt.Errorf("error opening device %04x:%04x: %w", vid, pid, err)
	}
	if dev == nil {
		return fmt.Errorf("printer %04x:%04x not found", vid, pid)
	}
	dev.Close()
	return nil
}
