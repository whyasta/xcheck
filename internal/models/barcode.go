package models

type Flag string
type Status string

const (
	BarcodeFlagValid   Flag = "VALID"
	BarcodeFlagUsed    Flag = "USED"
	BarcodeFlagExpired Flag = "EXPIRED"
)

const (
	BarcodeStatusIn   Status = "IN"
	BarcodeStatusOut  Status = "OUT"
	BarcodeStatusNull Status = ""
)

type Barcode struct {
	ID                int64            `gorm:"column:id; primary_key; not null" json:"id"`
	Barcode           string           `gorm:"column:barcode" json:"barcode" validate:"required"`
	Flag              Flag             `gorm:"column:flag;" json:"flag"`
	CurrentStatus     Status           `gorm:"column:current_status;" json:"current_status"`
	EventAssignmentID int64            `gorm:"column:event_assignment_id" json:"event_assignment_id"`
	EventAssignment   *EventAssignment `gorm:"foreignKey:id;references:event_assignment_id" json:"event_assignment"`
	CommonModel
}

type BarcodeAssignment struct {
	EventAssignmentID int64 `json:"event_assignment_id" validate:"required"`
	ImportId          int64 `json:"import_id" validate:"required"`
}
