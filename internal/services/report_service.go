package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
)

type ReportService struct {
	b repositories.BaseRepository
	r repositories.ReportRepository
}

func NewReportService(b repositories.BaseRepository, r repositories.ReportRepository) *ReportService {
	return &ReportService{b, r}
}

func (r *ReportService) ReportTrafficVisitor(eventID int64) (dto.TrafficVisitorSummary, error) {
	var data dto.TrafficVisitorSummary
	var bySession []dto.TrafficVisitorSession
	var byGate []dto.TrafficVisitorGate
	var byTicketType []dto.TrafficVisitorTicketType

	bySession, err := r.r.TrafficBySession(eventID)
	if err != nil {
		return data, err
	}

	byGate, err = r.r.TrafficByGate(eventID)
	if err != nil {
		return data, err
	}

	byTicketType, err = r.r.TrafficByTicketType(eventID)
	if err != nil {
		return data, err
	}

	data.Session = bySession
	data.Gate = byGate
	data.TicketType = byTicketType

	return data, nil
}

func (r *ReportService) ReportUniqueVisitor(eventID int64, ticketTypeIds []int64, gateIds []int64, sessionIds []int64) ([]dto.UniqueVisitorTicketType, error) {
	return r.r.UniqueByTicketType(eventID, ticketTypeIds, gateIds, sessionIds)
}

func (r *ReportService) ReportGateIn(eventID int64) ([]dto.GateInChart, error) {
	return r.r.GateIn(eventID)
}

func (r *ReportService) ReportRedemptionSummary(eventID int64, ticketTypeIds []int64, userIds []int64) ([]dto.RedemptionSummary, error) {
	return r.r.RedemptionSummary(eventID, ticketTypeIds, userIds)
}

func (r *ReportService) ReportRedemptionLog(eventID int64, orderBarcode string, orderID string) ([]models.Ticket, error) {
	return r.r.RedemptionLog(eventID, orderBarcode, orderID)
}
