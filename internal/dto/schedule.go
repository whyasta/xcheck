package dto

import "bigmind/xcheck-be/internal/models"

type ScheduleRequest struct {
	EventID      int64  `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	TicketTypeID int64  `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	SessionID    int64  `gorm:"column:session_id" mapstructure:"session_id" json:"session_id" validate:"required"`
	GateID       int64  `gorm:"column:gate_id" mapstructure:"gate_id" json:"gate_id" validate:"required"`
	UserID       *int64 `gorm:"column:user_id" mapstructure:"user_id" json:"user_id,omitempty"`
}

func (s *ScheduleRequest) ToEntity() *models.Schedule {
	return &models.Schedule{
		EventID:      s.EventID,
		TicketTypeID: s.TicketTypeID,
		SessionID:    s.SessionID,
		GateID:       s.GateID,
		UserID:       s.UserID,
	}
}
