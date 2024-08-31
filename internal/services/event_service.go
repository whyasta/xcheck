package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
)

type EventService struct {
	r repositories.EventRepository
}

func NewEventService(r repositories.EventRepository) *EventService {
	return &EventService{r}
}

func (s *EventService) CreateEvent(event *dto.EventRequest) (models.Event, error) {
	return s.r.Save(event)
}

func (s *EventService) CreateBulkEvent(events *[]dto.EventRequest) ([]models.Event, error) {
	return s.r.BulkSave(events)
}

func (s *EventService) UpdateEvent(id int64, event *map[string]interface{}) (models.Event, error) {
	// config.Logger.Infof("UpdateEvent: %+v", event)
	return s.r.Update(id, event)
}

func (s *EventService) GetAllEvents(pageParams *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error) {
	result, count, err := s.r.Paginate(pageParams, params)
	return result, count, err
}

func (s *EventService) GetFilteredEvents(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]dto.EventResponse, int64, error) {
	result, count, err := s.r.GetFiltered(pageParams, filters, sorts)

	rows := []dto.EventResponse{}
	for _, item := range result {
		// var gates []dto.EventGateResponse
		// for _, gate := range item.Gates {
		// 	gates = append(gates, dto.EventGateResponse{
		// 		ID:          gate.ID,
		// 		GateName:    gate.GateName,
		// 		EventID:     int64(item.ID),
		// 		TicketTypes: s.r.GateTicketTypes(int64(item.ID), gate.ID),
		// 	})
		// }

		rows = append(rows, dto.EventResponse{
			ID:          item.ID,
			EventName:   item.EventName,
			Status:      item.Status,
			StartDate:   item.StartDate.Format("2006-01-02"),
			EndDate:     item.EndDate.Format("2006-01-02"),
			TicketTypes: item.TicketTypes,
			Gates:       item.Gates,
			Sessions:    item.Sessions,
			EventSummary: dto.EventSummary{
				TotalBarcode:  0,
				TotalCheckIn:  0,
				TotalCheckOut: 0,
			},
		})
	}

	return rows, count, err
}

func (s *EventService) GetEventByID(uid int64) (dto.EventResponse, error) {
	res, err := s.r.FindByID(uid)
	if err != nil {
		return dto.EventResponse{}, err
	}

	// var gates []dto.EventGateResponse
	// for _, gate := range res.Gates {
	// 	gates = append(gates, dto.EventGateResponse{
	// 		ID:          gate.ID,
	// 		GateName:    gate.GateName,
	// 		EventID:     int64(res.ID),
	// 		TicketTypes: s.r.GateTicketTypes(int64(res.ID), gate.ID),
	// 	})
	// }
	gates := models.Gates(res.Gates)
	row := dto.EventResponse{
		ID:              res.ID,
		EventName:       res.EventName,
		Status:          res.Status,
		StartDate:       res.StartDate.Format("2006-01-02"),
		EndDate:         res.EndDate.Format("2006-01-02"),
		TicketTypes:     res.TicketTypes,
		Gates:           res.Gates,
		Sessions:        res.Sessions,
		EventSummary:    s.r.Summary(res.ID),
		LastSyncedAt:    res.LastSyncedAt,
		GateTicketTypes: s.r.GateTicketTypes(res.ID, gates.IDList()),
		// EventSummary: dto.EventSummary{
		// 	TotalBarcode:  0,
		// 	TotalCheckIn:  0,
		// 	TotalCheckOut: 0,
		// },
	}
	return row, nil
}

func (s *EventService) Delete(uid int64) (models.Event, error) {
	return s.r.Delete(uid)
}

func (s *EventService) Report(uid int64) (dto.EventReportResponse, error) {
	return s.r.Report(uid)
}
