package models

import "time"

type Schedule struct {
	ID               int64     `gorm:"column:id; primary_key; not null" json:"id"`
	ScheduleDateTime time.Time `gorm:"column:schedule_date_time" json:"schedule_date_time" validate:"required"`
	EventID          int64     `gorm:"column:event_id" json:"event_id" validate:"required"`
	CommonModel
}
