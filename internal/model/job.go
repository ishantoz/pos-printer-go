package model

import "time"

type BarcodeJob struct {
	ID             int       `json:"id"`
	VID            string    `json:"vid"`
	PID            string    `json:"pid"`
	SizeX          int       `json:"sizeX"`
	SizeY          int       `json:"sizeY"`
	Direction      int       `json:"direction"`
	TopText        string    `json:"topText"`
	BarcodeData    string    `json:"barcodeData"`
	PrintCount     int       `json:"printCount"`
	LabelGapLength int       `json:"labelGapLength"`
	LabelGapOffset int       `json:"labelGapOffset"`
	Status         string    `json:"status"`
	Attempts       int       `json:"attempts"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
