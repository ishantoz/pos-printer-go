package model

type LabelGap struct {
	Length int `json:"length"` // mm; 0 => auto-detect
	Offset int `json:"offset"`
}

type PrintBarcodeRequest struct {
	VID         string   `json:"vid"`
	PID         string   `json:"pid"`
	SizeX       int      `json:"sizeX"`
	SizeY       int      `json:"sizeY"`
	Direction   int      `json:"direction"`
	TopText     string   `json:"topText"`
	BarcodeData string   `json:"barcodeData"`
	PrintCount  int      `json:"printCount"`
	LabelGap    LabelGap `json:"labelGap"`
}
