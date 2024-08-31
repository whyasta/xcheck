package models

type GateAllocation struct {
	ID        int64    `gorm:"column:id; primary_key; not null" json:"id"`
	EventID   int64    `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	SessionID int64    `gorm:"column:session_id" mapstructure:"session_id" json:"session_id" validate:"required"`
	GateID    int64    `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
	Event     *Event   `gorm:"foreignKey:event_id;references:id" json:"event,omitempty"`
	Session   *Session `gorm:"foreignKey:session_id;references:id" json:"session,omitempty"`
	Gate      *Gate    `gorm:"foreignKey:gate_id;references:id" json:"gate,omitempty"`
	CommonModel
}
