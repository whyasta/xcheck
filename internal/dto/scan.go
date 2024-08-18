package dto

type ScanBarcode struct {
	EventID int64  `json:"event_id" validate:"required"`
	GateID  int64  `json:"gate_id" validate:"required"`
	Barcode string `json:"barcode" validate:"required"`
	Device  string `json:"device"`
}
