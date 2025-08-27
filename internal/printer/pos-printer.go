package printer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/gousb"
)

type usbWriter struct {
	ep *gousb.OutEndpoint
}

func (w *usbWriter) Write(p []byte) (int, error) {
	return w.ep.Write(p)
}

type PosPrinter struct {
	ctx *gousb.Context
}

func NewPosPrinter() *PosPrinter {
	return &PosPrinter{
		ctx: nil,
	}
}

func (p *PosPrinter) posPrinterContext(vidHexStr, pidHexStr string) (*gousb.Context, gousb.ID, gousb.ID, error) {
	vid64, err := strconv.ParseUint(vidHexStr, 0, 16)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("invalid Vendor ID %q: %w", vidHexStr, err)
	}
	pid64, err := strconv.ParseUint(pidHexStr, 0, 16)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("invalid Product ID %q: %w", pidHexStr, err)
	}

	_vid := gousb.ID(uint16(vid64))
	_pid := gousb.ID(uint16(pid64))

	ctx := gousb.NewContext()

	return ctx, _vid, _pid, nil
}

func (p *PosPrinter) OpenPosPrinter(vidHexStr, pidHexStr string) (*gousb.Device, error) {
	if p.ctx == nil {
		ctx, _, _, err := p.posPrinterContext(vidHexStr, pidHexStr)
		if err != nil {
			return nil, err
		}
		p.ctx = ctx
	}

	vid64, err := strconv.ParseUint(vidHexStr, 0, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid Vendor ID %q: %w", vidHexStr, err)
	}
	pid64, err := strconv.ParseUint(pidHexStr, 0, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid Product ID %q: %w", vidHexStr, err)
	}

	_vid := gousb.ID(uint16(vid64))
	_pid := gousb.ID(uint16(pid64))

	dev, err := p.ctx.OpenDeviceWithVIDPID(_vid, _pid)
	if err != nil {
		if p.ctx != nil {
			log.Printf("Failed to open device, resetting context and retrying: %v", err)
			p.ResetContext()

			// Create new context and retry
			ctx, _, _, err2 := p.posPrinterContext(vidHexStr, pidHexStr)
			if err2 != nil {
				return nil, fmt.Errorf("failed to create new context after reset: %w", err2)
			}
			p.ctx = ctx

			dev, err = p.ctx.OpenDeviceWithVIDPID(_vid, _pid)
			if err != nil {
				return nil, fmt.Errorf("error opening device %04x:%04x after retry: %w", _vid, _pid, err)
			}
		} else {
			return nil, fmt.Errorf("error opening device %04x:%04x: %w", _vid, _pid, err)
		}
	}
	return dev, nil
}

func (p *PosPrinter) CheckPrinter(vidHexStr, pidHexStr string) error {
	ctx, vid, pid, err := p.posPrinterContext(vidHexStr, pidHexStr)
	if err != nil {
		return err
	}
	defer ctx.Close()

	dev, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		log.Printf("error opening device %s:%s: %v", vidHexStr, pidHexStr, err)
		return err
	}
	if dev == nil {
		return fmt.Errorf("printer %s:%s not found", vidHexStr, pidHexStr)
	}
	defer dev.Close()
	return nil
}

func (p *PosPrinter) GetESCWriter(vidHexStr, pidHexStr string) (*usbWriter, *gousb.Device, error) {
	dev, err := p.OpenPosPrinter(vidHexStr, pidHexStr)
	if err != nil {
		return nil, nil, err
	}

	cfg, err := dev.Config(1)
	if err != nil {
		dev.Close()
		return nil, nil, err
	}

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		cfg.Close()
		dev.Close()
		return nil, nil, err
	}

	ep, err := intf.OutEndpoint(1)
	if err != nil {
		intf.Close()
		cfg.Close()
		dev.Close()
		return nil, nil, err
	}

	writer := &usbWriter{ep: ep}

	return writer, dev, nil
}

func (p *PosPrinter) Close() {
	if p.ctx != nil {
		p.ctx.Close()
		p.ctx = nil
	}
}

func (p *PosPrinter) Cleanup() {
	p.Close()
}

func (p *PosPrinter) IsReady() bool {
	return p.ctx != nil
}

func (p *PosPrinter) ResetContext() {
	if p.ctx != nil {
		p.ctx.Close()
		p.ctx = nil
	}
}
