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
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "DATA_NOT_FOUND", "UNKNOWN_ERROR", "INVALID_REQUEST", "UNAUTHORIZED", "SESSION_EXPIRED", "CHECKIN", "RE_CHECKIN", "CHECKOUT", "FAILED"}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Unknown Error", "Invalid Request", "Unauthorized", "SESSION_EXPIRED", "CHECKIN", "RE_CHECKIN", "CHECKOUT", "FAILED"}[r-1]
}
