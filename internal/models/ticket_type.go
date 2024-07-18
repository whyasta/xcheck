package models

type TicketType struct {
	ID             int64  `gorm:"column:id; primary_key; not null" mapstructure:"id" json:"id"`
	TicketTypeName string `gorm:"column:ticket_type_name" mapstructure:"ticket_type_name" json:"ticket_type_name" validate:"required,min=3,max=20"`
	EventID        int64  `gorm:"column:event_id" mapstructure:"event_id" json:"event_id" validate:"required"`
	CommonModel
}
