package models

import (
	"bigmind/xcheck-be/internal/constant"
	"time"
)

type LatestScan struct {
	Barcode       string                 `json:"barcode"`
	ScannedAt     time.Time              `json:"scanned_at"`
	GateID        int64                  `json:"gate_id"`
	GateName      string                 `json:"gate_name"`
	ScannedBy     int64                  `json:"scanned_by"`
	Device        string                 `json:"device"`
	Action        constant.BarcodeStatus `json:"action"`
	ScannedByName string                 `json:"scanned_by_name,omitempty"`
}

func (LatestScan) TableName() string {
	return "barcode_logs"
}

type Barcode struct {
	ID            int64                  `gorm:"column:id; primary_key; not null" json:"id"`
	Barcode       string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	Flag          constant.BarcodeFlag   `gorm:"column:flag;" json:"flag"`
	CurrentStatus constant.BarcodeStatus `gorm:"column:current_status;" json:"current_status"`
	EventID       int64                  `gorm:"column:event_id"  mapstructure:"event_id" json:"event_id" validate:"required"`
	TicketTypeID  int64                  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketType    *TicketType            `gorm:"foreignKey:id;references:ticket_type_id" json:"ticket_type"`
	Gates         []Gate                 `gorm:"many2many:barcode_gates;" json:"gates,omitempty"`
	Sessions      []Session              `gorm:"many2many:barcode_sessions;" json:"sessions,omitempty"`
	LatestScan    *LatestScan            `gorm:"foreignKey:barcode;references:barcode" json:"latest_scan"`
	CreatedAt     time.Time              `gorm:"column:created_at;column:created_at" json:"created_at"`
	// Sessions      []int64                `gorm:"serializer:json" mapstructure:"sessions" json:"sessions,omitempty"`
	// Gates         []int64                `gorm:"serializer:json" mapstructure:"gates" json:"gates,omitempty"`

	// GateAllocationID int64                  `gorm:"column:gateAllocation_id" json:"gateAllocation_id"`
	// GateAllocation   *GateAllocation        `gorm:"foreignKey:id;references:gate_allocation_id" json:"gateAllocation"`
}

type BarcodeAssignment struct {
	GateAllocationID int64 `json:"gateAllocation_id" validate:"required"`
	ImportID         int64 `json:"import_id" validate:"required"`
	TicketTypeID     int64 `json:"ticket_type_id" validate:"required"`
}

type BarcodeLog struct {
	ID           int64                  `json:"id,omitempty"`
	Barcode      string                 `json:"barcode"`
	EventID      int64                  `json:"event_id"`
	GateID       int64                  `json:"gate_id"`
	TicketTypeID int64                  `json:"ticket_type_id"`
	SessionID    int64                  `json:"session_id"`
	ScannedBy    int64                  `json:"scanned_by"`
	ScannedAt    time.Time              `json:"scanned_at"`
	Device       string                 `json:"device"`
	Action       constant.BarcodeStatus `json:"action"`
}
