package dto

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"time"
)

type BarcodeDownloadDto struct {
	EventID   int64 `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	SessionID int64 `gorm:"column:session_id" mapstructure:"session_id" json:"session_id" validate:"required"`
	GateID    int64 `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
}

type BarcodeUploadLogDto struct {
	Barcode   string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	ScannedAt time.Time              `gorm:"->:false;column:scanned_at" json:"scanned_at,omitempty"`
	GateID    int64                  `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
	EventID   int64                  `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	ScannedBy int64                  `gorm:"column:scanned_by" mapstructure:"scanned_by" json:"scanned_by" validate:"required"`
	Action    constant.BarcodeStatus `gorm:"column:action" mapstructure:"action" json:"action" validate:"required"`
}

type BarcodeUploadDto struct {
	Data []BarcodeUploadLogDto `gorm:"column:data" json:"data" validate:"required,dive"`
}

func (s *BarcodeUploadLogDto) ToEntity() *models.BarcodeLog {
	return &models.BarcodeLog{
		Barcode:   s.Barcode,
		ScannedAt: s.ScannedAt,
		ScannedBy: s.ScannedBy,
		Action:    s.Action,
		GateID:    s.GateID,
		EventID:   s.EventID,
	}
}
