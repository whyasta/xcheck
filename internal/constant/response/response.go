package response

type ResponseStatus int
type Headers int
type General int

// Constant Api
const (
	Success ResponseStatus = iota + 1
	DataNotFound
	UnknownError
	InvalidRequest
	Unauthorized
	SessionExpired
	Checkin
	ReCheckin
	Checkout
	Failed
	EC01
	EC02
	EC03
	EC04
	EC05
	EC11
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "UNKNOWN_ERROR", "INVALID_REQUEST",
		"UNAUTHORIZED", "SESSION_EXPIRED", "CHECKIN", "RE_CHECKIN", "CHECKOUT",
		"FAILED",
		"EC01", "EC02", "EC03",
		"EC04", "EC05", "EC11"}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Unknown Error", "Invalid Request", "Unauthorized",
		"SESSION_EXPIRED", "CHECKIN", "RE_CHECKIN", "CHECKOUT", "FAILED",
		"Barcode %s not found!",
		"Barcode found but ticket type not allowed",
		"Barcode %s not allowed to re-enter!",
		"Barcode found but session %s not match!",
		"Barcode found but gate %s not match",
		"Barcode %s not checked-in yet!"}[r-1]
}
