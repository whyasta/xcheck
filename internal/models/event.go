package models

import "time"

// swagger:model
type Event struct {
	ID        int64     `gorm:"column:id; primary_key; not null" json:"id"`
	EventName string    `gorm:"column:event_name" json:"event_name" validate:"required,min=5,max=20"`
	Status    int       `gorm:"column:status;default:0" json:"status"`
	StartDate time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   time.Time `gorm:"column:end_date" json:"end_date"`
	CommonModel
}

// swagger:model
type EventRequest struct {
	EventName string    `gorm:"column:event_name" json:"event_name" validate:"required,min=5,max=20"`
	Status    int       `gorm:"column:status;default:0" json:"status"`
	StartDate time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   time.Time `gorm:"column:end_date" json:"end_date"`
}

// swagger:parameters getEvent deleteEvent
type EventID struct {
	// In: path
	ID int `json:"id"`
}

// swagger:parameters createEvent
type EventCreateBodyParams struct {
	// required: true
	// in: body
	EventRequest *EventRequest `json:"EventRequest"`
}
