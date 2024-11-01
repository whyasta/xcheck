package repositories

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ReportRepository interface {
	TrafficBySession(eventID int64) ([]dto.TrafficVisitorSession, error)
	TrafficByGate(eventID int64) ([]dto.TrafficVisitorGate, error)
	TrafficByTicketType(eventID int64) ([]dto.TrafficVisitorTicketType, error)
	UniqueByTicketType(eventID int64, ticketTypeIds []int64, gateIds []int64, sessionIds []int64) ([]dto.UniqueVisitorTicketType, error)
	GateIn(eventID int64) ([]dto.GateInChart, error)

	RedemptionSummary(eventID int64, ticketTypeIds []int64, userIds []int64) ([]dto.RedemptionSummary, error)
	RedemptionLog(eventID int64, orderBarcode string, orderID string) ([]models.Ticket, error)
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
		Where("ticket_types.event_id = ? AND (barcode_logs.reason = '' OR barcode_logs.reason IS NULL)", eventID).
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

	// query := repo.db.Table("barcodes").
	// 	Debug().
	// 	Select("barcodes.ticket_type_id", "ticket_type_name",
	// 		"SUM(CASE WHEN current_status = 'IN' THEN 1 ELSE 0 END) as check_in_count",
	// 		"SUM(CASE WHEN current_status = 'OUT' THEN 1 ELSE 0 END) as check_out_count").
	// 	Joins("join ticket_types on ticket_types.id = barcodes.ticket_type_id").
	// 	Group("barcodes.ticket_type_id").
	// 	Where("barcodes.event_id = ?", eventID)

	sessionQuery := ""
	if len(ticketTypeIds) > 0 {
		var strNumbers []string
		for _, num := range ticketTypeIds {
			strNumbers = append(strNumbers, strconv.FormatInt(num, 10))
		}
		sessionQuery = sessionQuery + fmt.Sprintf(" and bl.ticket_type_id in (%s)", strings.Join(strNumbers, ","))
	}
	if len(gateIds) > 0 {
		var strNumbers []string
		for _, num := range gateIds {
			strNumbers = append(strNumbers, strconv.FormatInt(num, 10))
		}
		sessionQuery = sessionQuery + fmt.Sprintf(" and bl.gate_id in (%s)", strings.Join(strNumbers, ","))
	}
	if len(sessionIds) > 0 {
		var strNumbers []string
		for _, num := range sessionIds {
			strNumbers = append(strNumbers, strconv.FormatInt(num, 10))
		}
		sessionQuery = sessionQuery + fmt.Sprintf(" and bl.session_id in (%s)", strings.Join(strNumbers, ","))
	}

	checkInSelect := fmt.Sprintf("IFNULL((select COUNT(DISTINCT barcode) from barcode_logs bl where bl.action = 'IN' and bl.ticket_type_id = barcode_logs.ticket_type_id %s GROUP BY ticket_type_id), 0) as check_in_count", sessionQuery)
	checkOutSelect := fmt.Sprintf("IFNULL((select COUNT(DISTINCT barcode) from barcode_logs bl where bl.action = 'OUT' and bl.ticket_type_id = barcode_logs.ticket_type_id %s GROUP BY ticket_type_id), 0) as check_out_count", sessionQuery)

	query := repo.db.Table("barcode_logs").
		Debug().
		Select("barcode_logs.ticket_type_id", "ticket_type_name", checkInSelect, checkOutSelect).
		Joins("join ticket_types on ticket_types.id = barcode_logs.ticket_type_id").
		Group("barcode_logs.ticket_type_id").
		Where("(barcode_logs.reason = '' OR barcode_logs.reason IS NULL)").
		Where("barcode_logs.event_id = ?", eventID)

	if len(ticketTypeIds) > 0 {
		query = query.Where("barcode_logs.ticket_type_id in (?)", ticketTypeIds)
	}

	if len(gateIds) > 0 {
		query = query.Where("barcode_logs.gate_id in (?)", gateIds)
	}

	if len(sessionIds) > 0 {
		query = query.Where("barcode_logs.session_id in (?)", sessionIds)
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
		Where("barcode_logs.event_id = ? AND barcode_logs.action = ? AND (barcode_logs.reason = '' OR barcode_logs.reason IS NULL)", eventID, "IN").
		Group("DATE_FORMAT(scanned_at, '%Y-%m-%d %H.%i')").Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *reportRepository) RedemptionSummary(eventID int64, ticketTypeIds []int64, userIds []int64) ([]dto.RedemptionSummary, error) {
	var data = make([]dto.RedemptionSummary, 0)

	query := repo.db.Table("tickets").
		Debug().
		Select("ticket_type_id, ticket_type_name, count(order_barcode) as total").
		Where("tickets.event_id = ? AND tickets.status = ?", eventID, "REDEEMED")

	if len(ticketTypeIds) > 0 {
		query = query.Where("tickets.ticket_type_id in (?)", ticketTypeIds)
	}

	if len(userIds) > 0 {
		query = query.Where("tickets.redeemed_by in (?)", userIds)
	}

	err := query.Group("ticket_type_id, ticket_type_name").Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *reportRepository) RedemptionLog(eventID int64, orderBarcode string, orderID string) ([]models.Ticket, error) {
	var data = make([]models.Ticket, 0)

	query := repo.db.Table("tickets").
		Debug().
		Preload("User", func(tx2 *gorm.DB) *gorm.DB {
			return tx2.Omit("Password", "AuthUuids")
		}).
		Select("tickets.*").
		Where("tickets.event_id = ? AND tickets.status = ?", eventID, "REDEEMED")

	subWhereClause := ""
	if orderBarcode != "" {
		subWhereClause = subWhereClause + fmt.Sprintf("tickets.order_barcode = '%s'", orderBarcode)
	}

	if orderID != "" {
		if subWhereClause != "" {
			subWhereClause = subWhereClause + " OR "
		}
		subWhereClause = subWhereClause + fmt.Sprintf("tickets.order_id = '%s'", orderID)
	}

	if subWhereClause != "" {
		// subWhereClause = "(" + subWhereClause + ")"
		query = query.Where(subWhereClause)
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
