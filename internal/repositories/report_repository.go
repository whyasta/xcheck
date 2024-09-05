package repositories

import (
	"bigmind/xcheck-be/internal/dto"

	"gorm.io/gorm"
)

type ReportRepository interface {
	TrafficBySession(eventID int64) ([]dto.TrafficVisitorSession, error)
	TrafficByGate(eventID int64) ([]dto.TrafficVisitorGate, error)
	TrafficByTicketType(eventID int64) ([]dto.TrafficVisitorTicketType, error)
	UniqueByTicketType(eventID int64, ticketTypeIds []int64, gateIds []int64, sessionIds []int64) ([]dto.UniqueVisitorTicketType, error)
	GateIn(eventID int64) ([]dto.GateInChart, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *reportRepository {
	return &reportRepository{
		db: db,
	}
}

func (repo *reportRepository) TrafficBySession(eventID int64) ([]dto.TrafficVisitorSession, error) {
	var data = make([]dto.TrafficVisitorSession, 0)

	err := repo.db.Table("sessions").
		Select("sessions.id", "sessions.session_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.session_id = sessions.id").
		Where("sessions.event_id = ?", eventID).
		Group("sessions.id").
		Group("sessions.session_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) TrafficByGate(eventID int64) ([]dto.TrafficVisitorGate, error) {
	var data = make([]dto.TrafficVisitorGate, 0)

	err := repo.db.Table("gates").
		Select("gates.id", "gates.gate_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.gate_id = gates.id").
		Where("gates.event_id = ?", eventID).
		Group("gates.id").
		Group("gates.gate_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) TrafficByTicketType(eventID int64) ([]dto.TrafficVisitorTicketType, error) {
	var data = make([]dto.TrafficVisitorTicketType, 0)
	// err := repo.db.Table("barcode_logs").
	// 	Select("barcode_logs.ticket_type_id", "ticket_types.ticket_type_name",
	// 		"SUM(CASE WHEN action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
	// 		"SUM(CASE WHEN action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
	// 	Joins("join ticket_types on ticket_types.id = barcode_logs.ticket_type_id").
	// 	Where("barcode_logs.event_id = ?", eventID).
	// 	Group("barcode_logs.ticket_type_id").
	// 	Scan(&data).Error

	err := repo.db.Table("ticket_types").
		Select("ticket_types.id", "ticket_types.ticket_type_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.ticket_type_id = ticket_types.id").
		Where("ticket_types.event_id = ?", eventID).
		Group("ticket_types.id").
		Group("ticket_types.ticket_type_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) UniqueByTicketType(eventID int64, ticketTypeIds []int64, gateIds []int64, sessionIds []int64) ([]dto.UniqueVisitorTicketType, error) {
	var data = make([]dto.UniqueVisitorTicketType, 0)

	// subQuery := repo.db.Table("barcode_logs").
	// 	Debug().Select("barcode, action, barcode_logs.ticket_type_id, ticket_types.ticket_type_name").
	// 	Joins("join ticket_types on ticket_types.id = barcode_logs.ticket_type_id").
	// 	Where("barcode_logs.event_id = ?", eventID).
	// 	Group("barcode_logs.ticket_type_id").
	// 	Group("barcode_logs.barcode").
	// 	Group("barcode_logs.action").
	// 	Order("barcode_logs.barcode")

	// if len(ticketTypeIds) > 0 {
	// 	subQuery = subQuery.Where("ticket_type_id in (?)", ticketTypeIds)
	// }

	// err := repo.db.Table("(?) as u", subQuery).
	// 	Debug().
	// 	Select("ticket_type_id", "ticket_type_name",
	// 		"SUM(CASE WHEN action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
	// 		"SUM(CASE WHEN action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
	// 	Group("u.ticket_type_id").
	// 	Group("u.ticket_type_name").
	// 	Scan(&data).Error

	query := repo.db.Table("barcodes").
		Debug().
		Select("barcodes.ticket_type_id", "ticket_type_name",
			"SUM(CASE WHEN current_status = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN current_status = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("join ticket_types on ticket_types.id = barcodes.ticket_type_id").
		// Joins("join barcode_logs on barcode_logs.event_id = barcodes.event_id AND barcode_logs.barcode = barcodes.barcode").
		Group("barcodes.ticket_type_id").
		Where("barcodes.event_id = ?", eventID)

	if len(ticketTypeIds) > 0 {
		query = query.Where("barcodes.ticket_type_id in (?)", ticketTypeIds)
	}

	if len(gateIds) > 0 {
		mySubquery := repo.db.Table("barcode_logs").Select("barcode").
			Where("barcode_logs.event_id = barcodes.event_id AND barcode_logs.barcode = barcodes.barcode and  barcode_logs.gate_id in (?)", gateIds)
		query = query.Where("EXISTS (?)", mySubquery)
		//query = query.Where("barcode_logs.gate_id in (?)", gateIds)
	}

	if len(sessionIds) > 0 {
		mySubquery := repo.db.Table("barcode_logs").Select("barcode").
			Where("barcode_logs.event_id = barcodes.event_id AND barcode_logs.barcode = barcodes.barcode and  barcode_logs.session_id in (?)", sessionIds)
		query = query.Where("EXISTS (?)", mySubquery)
		//query = query.Where("barcode_logs.session_id in (?)", sessionIds)
	}

	err := query.Scan(&data).Error
	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) GateIn(eventID int64) ([]dto.GateInChart, error) {
	var data = make([]dto.GateInChart, 0)

	err := repo.db.Table("barcode_logs").
		Debug().
		Select("DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i') AS date_time", "COUNT(DISTINCT barcode) as total",
			"DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i') AS date_time", "COUNT(DISTINCT barcode) as unique_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as traffic_in_count").
		Where("barcode_logs.event_id = ?", eventID).
		Where("barcode_logs.action = ?", "IN").
		Group("DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i')").Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
