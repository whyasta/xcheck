package dto

type TicketRequest struct {
	OrderID         string `json:"order_id" validate:"required"`
	GenerateBarcode *bool  `json:"generate_barcode" validate:"required"`
}

type TicketTicketResponse struct {
	ID           int64  `json:"id"`
	OrderID      string `json:"order_id" validate:"required"`
	EventID      int64  `json:"event_id,omitempty"`
	TicketTypeID int64  `json:"ticket_type_id"`
}
