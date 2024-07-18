package dto

import (
	"bigmind/xcheck-be/internal/constant"
	"time"
)

type BarcodeLog struct {
	Barcode   string
	ScannedAt time.Time
	Action    constant.BarcodeStatus
}
