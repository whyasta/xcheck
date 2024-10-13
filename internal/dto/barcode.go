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
	ID           int64                  `json:"id,omitempty"`
	Barcode      string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	ScannedAt    time.Time              `gorm:"->:false;column:scanned_at" json:"scanned_at,omitempty"`
	GateID       int64                  `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
	EventID      int64                  `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	TicketTypeID int64                  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	SessionID    int64                  `gorm:"column:session_id" mapstructure:"session_id" json:"session_id" validate:"required"`
	ScannedBy    int64                  `gorm:"column:scanned_by" mapstructure:"scanned_by" json:"scanned_by" validate:"required"`
	Device       string                 `gorm:"column:device" mapstructure:"device" json:"device,omitempty"`
	Action       constant.BarcodeStatus `gorm:"column:action" mapstructure:"action" json:"action" validate:"required"`
}

type BarcodeUploadDataDto struct {
	ID            int64                  `gorm:"column:id; primary_key; not null" json:"id"`
	Barcode       string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	Flag          constant.BarcodeFlag   `gorm:"column:flag;" json:"flag"`
	CurrentStatus constant.BarcodeStatus `gorm:"column:current_status;" json:"current_status"`
	EventID       int64                  `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id" validate:"required"`
	TicketTypeID  int64                  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketType    *models.TicketType     `gorm:"foreignKey:id;references:ticket_type_id" json:"ticket_type"`
	Gates         *[]models.Gate         `gorm:"serializer:json" mapstructure:"gates" json:"gates,omitempty"`
	Sessions      *[]models.Session      `gorm:"serializer:json" mapstructure:"sessions" json:"sessions,omitempty"`
}

type BarcodeResponseDto struct {
	ID            int64                  `gorm:"column:id; primary_key; not null" json:"id"`
	Barcode       string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	Flag          constant.BarcodeFlag   `gorm:"column:flag;" json:"flag"`
	CurrentStatus constant.BarcodeStatus `gorm:"column:current_status;" json:"current_status"`
	EventID       int64                  `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id" validate:"required"`
	TicketTypeID  int64                  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketType    *models.TicketType     `gorm:"foreignKey:id;references:ticket_type_id" json:"ticket_type"`
	Gates         *[]models.Gate         `gorm:"serializer:json" mapstructure:"gates" json:"gates,omitempty"`
	Sessions      *[]models.Session      `gorm:"serializer:json" mapstructure:"sessions" json:"sessions,omitempty"`
}

type BarcodeUploadDto struct {
	Barcodes []BarcodeUploadDataDto `json:"barcodes" validate:"required,dive"`
	History  []BarcodeUploadLogDto  `json:"history" validate:"required,dive"`
}

type BarcodeLogResponseDto struct {
	ID            int64                  `json:"id,omitempty"`
	Barcode       string                 `json:"barcode"`
	EventID       int64                  `json:"event_id"`
	GateID        int64                  `json:"gate_id"`
	TicketTypeID  int64                  `json:"ticket_type_id"`
	SessionID     int64                  `json:"session_id"`
	ScannedBy     int64                  `json:"scanned_by"`
	ScannedAt     time.Time              `json:"scanned_at"`
	Device        string                 `json:"device"`
	Action        constant.BarcodeStatus `json:"action"`
	CurrentStatus constant.BarcodeStatus `json:"current_status"`
	Flag          constant.BarcodeFlag   `json:"flag"`
}

type BarcodeUpdateRequest struct {
	Gates    []int64 `mapstructure:"gates" json:"gates,omitempty" validate:"required"`
	Sessions []int64 `mapstructure:"sessions" json:"sessions,omitempty" validate:"required"`
}

type BarcodeUpdateByTicketTypeRequest struct {
	TicketTypeID int64   `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	Gates        []int64 `mapstructure:"gates" json:"gates,omitempty" validate:"required"`
	Sessions     []int64 `mapstructure:"sessions" json:"sessions,omitempty" validate:"required"`
}

func (s *BarcodeUploadDataDto) ToEntity() *models.Barcode {
	return &models.Barcode{
		Barcode:       s.Barcode,
		Flag:          s.Flag,
		CurrentStatus: s.CurrentStatus,
		EventID:       s.EventID,
		TicketTypeID:  s.TicketTypeID,
		Gates:         *s.Gates,
		Sessions:      *s.Sessions,
	}
}

func (s *BarcodeUploadLogDto) ToEntity() *models.BarcodeLog {
	return &models.BarcodeLog{
		ID:           s.ID,
		Barcode:      s.Barcode,
		ScannedAt:    s.ScannedAt,
		ScannedBy:    s.ScannedBy,
		Action:       s.Action,
		GateID:       s.GateID,
		EventID:      s.EventID,
		TicketTypeID: s.TicketTypeID,
		Device:       s.Device,
		SessionID:    s.SessionID,
	}
}
