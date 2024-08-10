package constant

type BarcodeFlag string
type BarcodeStatus string
type ImportStatus string

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
