package models

import "gorm.io/datatypes"

type Redeem struct {
	ID           int64          `gorm:"column:id; primary_key; not null" json:"id"`
	OrderID      string         `gorm:"column:order_id" json:"order_id" validate:"required"`
	EventID      int64          `gorm:"column:event_id" mapstructure:"event_id" json:"event_id,omitempty" validate:"required"`
	TicketTypeID int64          `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketType   *TicketType    `gorm:"foreignKey:id;references:ticket_type_id" json:"ticket_type"`
	Attributes   datatypes.JSON `gorm:"column:attributes" json:"attributes"`
	CommonModel
}

func (Redeem) TableName() string {
	return "redemption_tickets"
}
