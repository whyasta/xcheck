package dto

import (
	"bigmind/xcheck-be/internal/models"
	"time"
)

// swagger:model
type EventRequest struct {
	ID        int64     `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string    `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int       `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty" validate:"required"`
	StartDate time.Time `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty" validate:"required"`
	EndDate   time.Time `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty" validate:"required"`
}

type EventUpdateDto struct {
	ID        int64     `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string    `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int       `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty"`
	StartDate time.Time `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty"`
	EndDate   time.Time `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty"`
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

func (s *EventRequest) ToEntity() *models.Event {
	return &models.Event{
		ID:        s.ID,
		EventName: s.EventName,
		Status:    s.Status,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
	}
}
