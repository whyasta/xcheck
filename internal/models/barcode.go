package models

import "bigmind/xcheck-be/internal/constant"

type Barcode struct {
	ID            int64                  `gorm:"column:id; primary_key; not null" json:"id"`
	Barcode       string                 `gorm:"column:barcode" json:"barcode" validate:"required"`
	Flag          constant.BarcodeFlag   `gorm:"column:flag;" json:"flag"`
	CurrentStatus constant.BarcodeStatus `gorm:"column:current_status;" json:"current_status"`
	ScheduleID    int64                  `gorm:"column:schedule_id" json:"schedule_id"`
	Schedule      *Schedule              `gorm:"foreignKey:id;references:schedule_id" json:"schedule"`
	CommonModel
}

type BarcodeAssignment struct {
	ScheduleID int64 `json:"schedule_id" validate:"required"`
	ImportId   int64 `json:"import_id" validate:"required"`
}
