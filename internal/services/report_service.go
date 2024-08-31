package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/repositories"
)

type ReportService struct {
	b repositories.BaseRepository
	r repositories.ReportRepository
}

func NewReportService(b repositories.BaseRepository, r repositories.ReportRepository) *ReportService {
	return &ReportService{b, r}
}

func (r *ReportService) ReportTrafficVisitor(eventId int64) (dto.TrafficVisitorSummary, error) {
	var data dto.TrafficVisitorSummary
	var bySession []dto.TrafficVisitorSession
	var byGate []dto.TrafficVisitorGate
	var byTicketType []dto.TrafficVisitorTicketType

	bySession, err := r.r.TrafficBySession(eventId)
	if err != nil {
		return data, err
	}

	byGate, err = r.r.TrafficByGate(eventId)
	if err != nil {
		return data, err
	}

	byTicketType, err = r.r.TrafficByTicketType(eventId)
	if err != nil {
		return data, err
	}

	data.Session = bySession
	data.Gate = byGate
	data.TicketType = byTicketType

	return data, nil
}

func (r *ReportService) ReportUniqueVisitor(eventId int64, ticketTypeIds []int64, gateIds []int64, sessionIds []int64) ([]dto.UniqueVisitorTicketType, error) {
	return r.r.UniqueByTicketType(eventId, ticketTypeIds, gateIds, sessionIds)
}

func (r *ReportService) ReportGateIn(eventId int64) ([]dto.GateInChart, error) {
	return r.r.GateIn(eventId)
}
