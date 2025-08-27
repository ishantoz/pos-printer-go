package main

import (
	"image"
	"image/color"
	"log"

	"pos-printer/internal/lib"

	"github.com/gen2brain/go-fitz"
	"golang.org/x/image/draw"
)

// ESC/POS commands
var (
	ESC_INIT    = []byte{0x1B, 0x40}       // ESC @ -> initialize
	ESC_ALIGN_L = []byte{0x1B, 0x61, 0x00} // left align
	CUT_FULL    = []byte{0x1D, 0x56, 0x00} // full cut
)

func ESC_FEED_N(n byte) []byte {
	return []byte{0x1B, 0x64, n} // ESC d n -> feed n lines
}

// ResizeToWidth scales src image to target width, preserving aspect ratio
func ResizeToWidth(src image.Image, targetWidth int) image.Image {
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	scale := float64(targetWidth) / float64(srcWidth)
	newHeight := int(float64(srcHeight) * scale)

	dst := image.NewGray(image.Rect(0, 0, targetWidth, newHeight))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

// ApplyThreshold converts to black/white based on threshold (0-255)
func ApplyThreshold(img image.Image, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if c.Y < threshold {
				gray.SetGray(x, y, color.Gray{Y: 0})
			} else {
				gray.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}
	return gray
}

// Convert image.Image to ESC/POS raster format (GS v 0)
func ImageToRaster(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	data := []byte{0x1D, 0x76, 0x30, 0x00}                  // GS v 0, mode 0
	data = append(data, byte(width/8), 0x00)                // width in bytes
	data = append(data, byte(height&0xFF), byte(height>>8)) // height LSB/MSB

	for y := 0; y < height; y++ {
		for xByte := 0; xByte < width/8; xByte++ {
			var b byte
			for bit := 0; bit < 8; bit++ {
				x := xByte*8 + bit
				c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
				if c.Y < 128 {
					b |= (1 << uint(7-bit))
				}
			}
			data = append(data, b)
		}
	}
	return data
}

// PrintPDFScaled prints a PDF file to a thermal printer
func PrintPDFScaled(pdfPath string, writer interface {
	Write([]byte) (int, error)
	Close() error
}, printerWidth int, threshold uint8) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return err
	}
	defer doc.Close()

	for i := 0; i < doc.NumPage(); i++ {
		img, err := doc.Image(i)
		if err != nil {
			log.Printf("Failed to render page %d: %v", i, err)
			continue
		}

		// Resize & threshold
		resized := ResizeToWidth(img, printerWidth)
		bw := ApplyThreshold(resized, threshold)

		// Send ESC/POS commands
		writer.Write(ESC_INIT)
		writer.Write(ESC_ALIGN_L)

		// Send image raster data
		raster := ImageToRaster(bw)
		writer.Write(raster)

		// Feed paper and cut
		writer.Write(ESC_FEED_N(5))
		writer.Write(CUT_FULL)
	}

	return nil
}

func main() {
	pdfFile := "example.pdf" // your PDF file path
	vid := "0x0fe6"          // printer Vendor ID
	pid := "0x811e"          // printer Product ID
	printerWidth := 576      // 80mm thermal printer typical width
	threshold := uint8(130)  // black/white threshold

	// Create HID printer connection directly
	writer, err := lib.NewHIDWriter(vid, pid)
	if err != nil {
		log.Fatalf("Failed to open printer: %v", err)
	}
	defer writer.Close()

	log.Println("Printer connected successfully!")

	// Print PDF
	if err := PrintPDFScaled(pdfFile, writer, printerWidth, threshold); err != nil {
		log.Fatalf("Failed to print PDF: %v", err)
	}

	log.Println("PDF printed successfully!")
}
