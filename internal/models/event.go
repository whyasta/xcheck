package models

import (
	"time"

	"gorm.io/datatypes"
)

type Event struct {
	ID           int64          `gorm:"column:id; primary_key; not null" json:"id"`
	EventName    string         `gorm:"column:event_name" json:"event_name" validate:"required,min=5,max=100"`
	Status       int            `gorm:"column:status;default:0" json:"status"`
	StartDate    time.Time      `gorm:"column:start_date" json:"start_date"`
	EndDate      time.Time      `gorm:"column:end_date" json:"end_date"`
	LastSyncedAt *time.Time     `gorm:"column:last_synced_at" json:"last_synced_at"`
	RedeemConfig datatypes.JSON `gorm:"column:redeem_config" json:"redeem_config"`
	TicketTypes  []TicketType   `gorm:"foreignKey:event_id;references:id" json:"ticket_types"`
	Gates        []Gate         `gorm:"foreignKey:event_id;references:id" json:"gates"`
	Sessions     []Session      `gorm:"foreignKey:event_id;references:id" json:"sessions"`
	CommonModel
}

func (Event) TableName() string {
	return "events"
}
