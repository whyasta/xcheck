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
	ScheduleID    int64                  `gorm:"column:schedule_id" json:"schedule_id"`
	TicketTypeID  int64                  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	Schedule      *Schedule              `gorm:"foreignKey:id;references:schedule_id" json:"schedule"`
	CommonModel
}

type BarcodeAssignment struct {
	ScheduleID   int64 `json:"schedule_id" validate:"required"`
	ImportId     int64 `json:"import_id" validate:"required"`
	TicketTypeID int64 `json:"ticket_type_id" validate:"required"`
}

type BarcodeLog struct {
	EventID   int64
	Barcode   string
	ScannedAt time.Time
	GateID    int64
	ScannedBy int64
	Action    constant.BarcodeStatus
}
