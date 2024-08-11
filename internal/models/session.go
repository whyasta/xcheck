package models

import "time"

type Session struct {
	ID           int64     `gorm:"column:id; primary_key; not null" json:"id"`
	EventID      int64     `gorm:"column:event_id" mapstructure:"event_id" json:"event_id,omitempty" validate:"required"`
	Sessioname   string    `gorm:"column:session_name" mapstructure:"session_name" json:"session_name" validate:"required,min=5,max=100"`
	SessionStart time.Time `gorm:"column:session_start" mapstructure:"session_start" json:"session_start"`
	SessionEnd   time.Time `gorm:"column:session_end" mapstructure:"session_end" json:"session_end"`
	CommonModel
}
