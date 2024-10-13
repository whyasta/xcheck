package constant

type BarcodeFlag string
type BarcodeStatus string
type ImportStatus string
type TicketStatus string

const (
	BarcodeFlagValid   BarcodeFlag = "VALID"
	BarcodeFlagUsed    BarcodeFlag = "USED"
	BarcodeFlagExpired BarcodeFlag = "EXPIRED"
)

const (
	BarcodeStatusIn   BarcodeStatus = "IN"
	BarcodeStatusOut  BarcodeStatus = "OUT"
	BarcodeStatusNull BarcodeStatus = ""
)

const (
	ImportStatusPending    ImportStatus = "PENDING"
	ImportStatusProcessing ImportStatus = "PROCESSING"
	ImportStatusCompleted  ImportStatus = "COMPLETED"
	ImportStatusAssigned   ImportStatus = "ASSIGNED"
	ImportStatusFailed     ImportStatus = "FAILED"
)

const (
	TicketStatusPending  TicketStatus = "PENDING"
	TicketStatusRedeemed TicketStatus = "REDEEMED"
	TicketStatusExpired  TicketStatus = "EXPIRED"
	TicketStatusCanceled TicketStatus = "CANCELED"
)
