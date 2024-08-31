package dto

import (
	"bigmind/xcheck-be/internal/models"
	"log"
	"time"
)

type EventRequest struct {
	ID        int64  `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int    `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty" validate:"required"`
	StartDate string `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty" validate:"required,date"`
	EndDate   string `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty" validate:"required,date"`
}

type EventGateTicketTypeResponse struct {
	GateID         int64  `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
	GateName       string `gorm:"column:gate_name"  mapstructure:"gate_name" json:"gate_name" validate:"required,min=3,max=20"`
	TicketTypeID   int64  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketTypeName string `gorm:"column:ticket_type_name" mapstructure:"ticket_type_name" json:"ticket_type_name" validate:"required,min=3,max=20"`
}

type EventResponse struct {
	ID              int64                         `json:"id"`
	EventName       string                        `json:"event_name" validate:"required,min=5,max=100"`
	Status          int                           `json:"status"`
	StartDate       string                        `json:"start_date"`
	EndDate         string                        `json:"end_date"`
	LastSyncedAt    *time.Time                    `json:"last_synced_at"`
	TicketTypes     []models.TicketType           `json:"ticket_types"`
	Gates           []models.Gate                 `json:"gates"`
	Sessions        []models.Session              `json:"sessions"`
	GateTicketTypes []EventGateTicketTypeResponse `json:"gate_ticket_types"`
	EventSummary    `json:"summary"`
}

type EventSummary struct {
	TotalBarcode    int64 `json:"total_barcode"`
	OngoingCheckIn  int64 `json:"ongoing_check_in"`
	OngoingCheckOut int64 `json:"ongoing_check_out"`
	TotalCheckIn    int64 `json:"total_check_in"`
	TotalCheckOut   int64 `json:"total_check_out"`
	// TotalTicketType []map[string]interface{} `json:"total_ticket_type" gorm:"serializer:json"`
}

type EventUpdateDto struct {
	ID        int64  `gorm:"column:id" mapstructure:"id" json:"id,omitempty"`
	EventName string `gorm:"column:event_name" mapstructure:"event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status    int    `gorm:"column:status;default:0" mapstructure:"status" json:"status,omitempty"`
	StartDate string `gorm:"column:start_date" mapstructure:"start_date" json:"start_date,omitempty" validate:"date"`
	EndDate   string `gorm:"column:end_date" mapstructure:"end_date" json:"end_date,omitempty" validate:"date"`
}

type EventID struct {
	// In: path
	ID int `json:"id"`
}

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

type EventReportResponse struct {
	ID           int64      `json:"id"`
	EventName    string     `json:"event_name" validate:"required,min=5,max=100"`
	Status       int        `json:"status"`
	StartDate    string     `json:"start_date"`
	EndDate      string     `json:"end_date"`
	LastSyncedAt *time.Time `json:"last_synced_at"`
	// TicketTypes  []models.TicketType `gorm:"foreignKey:event_id;references:id" json:"ticket_types"`
	// Gates        []models.Gate       `gorm:"foreignKey:event_id;references:id" json:"gates"`
	// Sessions     []models.Session    `gorm:"foreignKey:event_id;references:id" json:"sessions"`
	EventSummary `json:"summary"`
}
