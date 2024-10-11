package models

import "gorm.io/datatypes"

type Ticket struct {
	ID             int64          `gorm:"column:id; primary_key; not null" json:"id"`
	ImportID       string         `gorm:"column:import_id" json:"import_id"`
	OrderBarcode   string         `gorm:"column:order_barcode" json:"order_barcode" validate:"required"`
	OrderID        string         `gorm:"column:order_id" json:"order_id" validate:"required"`
	EventID        int64          `gorm:"column:event_id" mapstructure:"event_id" json:"event_id,omitempty" validate:"required"`
	TicketTypeID   *int64         `gorm:"column:ticket_type_id" mapstructure:"ticket_type_id" json:"ticket_type_id" validate:"required"`
	TicketType     *TicketType    `gorm:"foreignKey:id;references:ticket_type_id" json:"-"`
	TicketTypeName string         `gorm:"column:ticket_type_name" mapstructure:"ticket_type_name" json:"ticket_type_name" validate:"required,min=3,max=100"`
	Attributes     datatypes.JSON `gorm:"column:attributes" json:"attributes"`
	Name           string         `gorm:"column:name" json:"name"`
	Email          string         `gorm:"column:email" json:"email"`
	PhoneNumber    string         `gorm:"column:phone_number" json:"phone_number"`
	Note           string         `gorm:"column:note" json:"note"`
	AssignStatus   int            `gorm:"column:assign_status" json:"assign_status"`
	Quantity       int            `gorm:"column:quantity" json:"quantity"`
	CommonModel
}

func (Ticket) TableName() string {
	return "tickets"
}
