package dto

import (
	"bigmind/xcheck-be/internal/constant"
	"time"
)

type BarcodeLog struct {
	Barcode   string
	ScannedAt time.Time
	ScannedBy int64
	Action    constant.BarcodeStatus
}

type BarcodeDownloadDto struct {
	EventID   int64 `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	SessionID int64 `gorm:"column:session_id" mapstructure:"session_id" json:"session_id" validate:"required"`
	GateID    int64 `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
}
