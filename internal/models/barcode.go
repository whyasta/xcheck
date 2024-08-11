package models

import (
	"bigmind/xcheck-be/internal/constant"
	"time"
)

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
	// Sessions      []int64                `gorm:"serializer:json" mapstructure:"sessions" json:"sessions,omitempty"`
	// Gates         []int64                `gorm:"serializer:json" mapstructure:"gates" json:"gates,omitempty"`

	// GateAllocationID int64                  `gorm:"column:gateAllocation_id" json:"gateAllocation_id"`
	// GateAllocation   *GateAllocation        `gorm:"foreignKey:id;references:gate_allocation_id" json:"gateAllocation"`
	CommonModel
}

type BarcodeAssignment struct {
	GateAllocationID int64 `json:"gateAllocation_id" validate:"required"`
	ImportId         int64 `json:"import_id" validate:"required"`
	TicketTypeID     int64 `json:"ticket_type_id" validate:"required"`
}

type BarcodeLog struct {
	EventID   int64
	Barcode   string
	ScannedAt time.Time
	GateID    int64
	ScannedBy int64
	Action    constant.BarcodeStatus
}
