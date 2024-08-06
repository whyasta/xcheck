package dto

import (
	"bigmind/xcheck-be/internal/models"
	"log"
	"time"
)

// swagger:model
type EventRequest struct {
	ID        int64  `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int    `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty" validate:"required"`
	StartDate string `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty" validate:"required,date"`
	EndDate   string `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty" validate:"required,date"`
}

type EventResponse struct {
	ID           int64               `json:"id"`
	EventName    string              `json:"event_name" validate:"required,min=5,max=100"`
	Status       int                 `json:"status"`
	StartDate    string              `json:"start_date"`
	EndDate      string              `json:"end_date"`
	TicketTypes  []models.TicketType `json:"ticket_types"`
	Gates        []models.Gate       `json:"gates"`
	Sessions     []models.Session    `json:"sessions"`
	EventSummary `json:"summary"`
}

type EventSummary struct {
	TotalBarcode  int64 `json:"total_barcode"`
	TotalCheckIn  int64 `json:"total_check_in"`
	TotalCheckOut int64 `json:"total_check_out"`
}

type EventUpdateDto struct {
	ID        int64  `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int    `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty"`
	StartDate string `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty" validate:"date"`
	EndDate   string `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty" validate:"date"`
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
	start, _ := time.Parse("2006-01-02", s.StartDate)
	end, _ := time.Parse("2006-01-02", s.EndDate)

	log.Println(s.StartDate, start)
	// config.Logger.Infof("ToEntity: %+v", s)
	return &models.Event{
		ID:        s.ID,
		EventName: s.EventName,
		Status:    s.Status,
		StartDate: start,
		EndDate:   end,
	}
}
