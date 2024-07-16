package models

type Gate struct {
	ID       int64  `gorm:"column:id; primary_key; not null" json:"id"`
	GateName string `gorm:"column:gate_name" json:"gate_name" validate:"required,min=3,max=20"`
	EventID  int64  `gorm:"column:event_id" json:"event_id" validate:"required"`
	CommonModel
}
