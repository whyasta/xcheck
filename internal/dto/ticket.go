package dto

type TicketCheckRequest struct {
	OrderBarcode string `json:"order_barcode" validate:"required"`
}

type TicketRedeemDataRequest struct {
	OrderID          string `json:"order_id" validate:"required"`
	AssociateBarcode string `json:"associate_barcode" validate:"required"`
}

type TicketRedeemRequest struct {
	GenerateBarcode bool                      `json:"generate_barcode" validate:"boolean"`
	Photo           *string                   `json:"photo" validate:"base64"`
	Note            *string                   `json:"note"`
	Data            []TicketRedeemDataRequest `json:"data" validate:"required,dive"`
}
