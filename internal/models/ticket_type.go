package models

type TicketType struct {
	ID             int64  `gorm:"column:id; primary_key; not null" json:"id"`
	TicketTypeName string `gorm:"column:ticket_type_name" json:"ticket_type_name" validate:"required,min=3,max=20"`
	CommonModel
}
