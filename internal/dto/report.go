package dto

type TrafficVisitorSession struct {
	SessionID     int64  `json:"session_id"`
	SessionName   string `json:"session_name"`
	CheckInCount  int64  `json:"check_in_count"`
	CheckOutCount int64  `json:"check_out_count"`
}

type TrafficVisitorGate struct {
	GateID        int64  `json:"gate_id"`
	GateName      string `json:"gate_name"`
	CheckInCount  int64  `json:"check_in_count"`
	CheckOutCount int64  `json:"check_out_count"`
}

type TrafficVisitorTicketType struct {
	TicketTypeID   int64  `json:"ticket_type_id"`
	TicketTypeName string `json:"ticket_type_name"`
	CheckInCount   int64  `json:"check_in_count"`
	CheckOutCount  int64  `json:"check_out_count"`
}

type UniqueVisitorTicketType struct {
	TicketTypeID   int64  `json:"ticket_type_id"`
	TicketTypeName string `json:"ticket_type_name"`
	CheckInCount   int64  `json:"check_in_count"`
	CheckOutCount  int64  `json:"check_out_count"`
}

type TrafficVisitorSummary struct {
	Session    []TrafficVisitorSession    `json:"session" gorm:"serializer:json"`
	Gate       []TrafficVisitorGate       `json:"gate" gorm:"serializer:json"`
	TicketType []TrafficVisitorTicketType `json:"ticket_type" gorm:"serializer:json"`
}

type GateInChart struct {
	DateTime string `json:"date_time"`
	Total    int64  `json:"total"`
}
