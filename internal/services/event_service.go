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
		rows = append(rows, dto.EventResponse{
			ID:          item.ID,
			EventName:   item.EventName,
			Status:      item.Status,
			StartDate:   item.StartDate,
			EndDate:     item.EndDate,
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

	row := dto.EventResponse{
		ID:           res.ID,
		EventName:    res.EventName,
		Status:       res.Status,
		StartDate:    res.StartDate,
		EndDate:      res.EndDate,
		TicketTypes:  res.TicketTypes,
		Gates:        res.Gates,
		Sessions:     res.Sessions,
		EventSummary: s.r.Summary(res.ID),
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
