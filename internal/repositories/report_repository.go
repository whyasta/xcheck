package repositories

import (
	"bigmind/xcheck-be/internal/dto"

	"gorm.io/gorm"
)

type ReportRepository interface {
	TrafficBySession(eventId int64) ([]dto.TrafficVisitorSession, error)
	TrafficByGate(eventId int64) ([]dto.TrafficVisitorGate, error)
	TrafficByTicketType(eventId int64) ([]dto.TrafficVisitorTicketType, error)
	UniqueByTicketType(eventId int64, ticketTypeIds []int64) ([]dto.UniqueVisitorTicketType, error)
	GateIn(eventId int64) ([]dto.GateInChart, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *reportRepository {
	return &reportRepository{
		db: db,
	}
}

func (repo *reportRepository) TrafficBySession(eventId int64) ([]dto.TrafficVisitorSession, error) {
	var data []dto.TrafficVisitorSession

	err := repo.db.Table("sessions").
		Select("sessions.id", "sessions.session_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.session_id = sessions.id").
		Where("sessions.event_id = ?", eventId).
		Group("sessions.id").
		Group("sessions.session_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) TrafficByGate(eventId int64) ([]dto.TrafficVisitorGate, error) {
	var data []dto.TrafficVisitorGate

	err := repo.db.Table("gates").
		Select("gates.id", "gates.gate_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.gate_id = gates.id").
		Where("gates.event_id = ?", eventId).
		Group("gates.id").
		Group("gates.gate_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) TrafficByTicketType(eventId int64) ([]dto.TrafficVisitorTicketType, error) {
	var data []dto.TrafficVisitorTicketType
	// err := repo.db.Table("barcode_logs").
	// 	Select("barcode_logs.ticket_type_id", "ticket_types.ticket_type_name",
	// 		"SUM(CASE WHEN action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
	// 		"SUM(CASE WHEN action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
	// 	Joins("join ticket_types on ticket_types.id = barcode_logs.ticket_type_id").
	// 	Where("barcode_logs.event_id = ?", eventId).
	// 	Group("barcode_logs.ticket_type_id").
	// 	Scan(&data).Error

	err := repo.db.Table("ticket_types").
		Select("ticket_types.id", "ticket_types.ticket_type_name",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("left join barcode_logs on barcode_logs.ticket_type_id = ticket_types.id").
		Where("ticket_types.event_id = ?", eventId).
		Group("ticket_types.id").
		Group("ticket_types.ticket_type_name").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) UniqueByTicketType(eventId int64, ticketTypeIds []int64) ([]dto.UniqueVisitorTicketType, error) {
	var data []dto.UniqueVisitorTicketType

	// subQuery := repo.db.Table("barcode_logs").
	// 	Debug().Select("barcode, action, barcode_logs.ticket_type_id, ticket_types.ticket_type_name").
	// 	Joins("join ticket_types on ticket_types.id = barcode_logs.ticket_type_id").
	// 	Where("barcode_logs.event_id = ?", eventId).
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
		Select("ticket_type_id", "ticket_type_name",
			"SUM(CASE WHEN current_status = 'IN' THEN 1 ELSE 0 END) as check_in_count",
			"SUM(CASE WHEN current_status = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
		Joins("join ticket_types on ticket_types.id = barcodes.ticket_type_id").
		Group("barcodes.ticket_type_id").
		Where("barcodes.event_id = ?", eventId)

	if len(ticketTypeIds) > 0 {
		query = query.Where("ticket_type_id in (?)", ticketTypeIds)
	}

	err := query.Scan(&data).Error
	if err != nil {
		return nil, err
	}

	return data, err
}

func (repo *reportRepository) GateIn(eventId int64) ([]dto.GateInChart, error) {
	var data []dto.GateInChart

	err := repo.db.Table("barcode_logs").
		Debug().
		Select("DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i') AS date_time", "COUNT(DISTINCT barcode) as total",
			"DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i') AS date_time", "COUNT(DISTINCT barcode) as unique_in_count",
			"SUM(CASE WHEN barcode_logs.action = 'IN' THEN 1 ELSE 0 END) as traffic_in_count").
		Where("barcode_logs.event_id = ?", eventId).
		Where("barcode_logs.action = ?", "IN").
		Group("DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i')").Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
