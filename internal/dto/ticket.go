package dto

type TicketCheckRequest struct {
	OrderBarcode string `json:"order_barcode" validate:"required"`
}

type TicketRedeemDataRequest struct {
	// OrderID          int64  `json:"order_id" validate:"required,numeric"`
	// OrderBarcode     string `json:"order_barcode" validate:"required"`
	ID               int64  `json:"id" validate:"required,numeric"`
	AssociateBarcode string `json:"associate_barcode"`
}

type TicketRedeemRequest struct {
	// GenerateBarcode bool                      `json:"generate_barcode" validate:"boolean"`
	// Photo *string                   `form:"photo" validate:"base64"`
	Note *string `form:"note" json:"note"`
	// Data []string `form:"data[].*"`
	Data []TicketRedeemDataRequest `json:"data" validate:"required,dive"`
}
