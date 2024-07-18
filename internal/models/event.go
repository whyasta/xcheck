package models

import "time"

// swagger:model
type Event struct {
	ID          int64         `gorm:"column:id; primary_key; not null" json:"id"`
	EventName   string        `gorm:"column:event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status      int           `gorm:"column:status;default:0" json:"status"`
	StartDate   time.Time     `gorm:"column:start_date" json:"start_date"`
	EndDate     time.Time     `gorm:"column:end_date" json:"end_date"`
	TicketTypes []*TicketType `gorm:"foreignKey:event_id;references:id" json:"ticket_types"`
	Gates       []*Gate       `gorm:"foreignKey:event_id;references:id" json:"gates"`
	Sessions    []*Session    `gorm:"foreignKey:event_id;references:id" json:"sessions"`
	CommonModel
}

func (Event) TableName() string {
	return "events"
}

// swagger:model
type EventRequest struct {
	ID        int       `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
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
