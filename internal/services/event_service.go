package services

import (
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

func (s *EventService) CreateEvent(role *models.Event) (models.Event, error) {
	return s.r.Save(role)
}

func (s *EventService) GetAllEvents(pageParams *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error) {
	result, count, err := s.r.Paginate(pageParams, params)
	return result, count, err
}

func (s *EventService) GetEventByID(uid int64) (models.Event, error) {
	return s.r.FindByID(uid)
}

func (s *EventService) Delete(uid int64) (models.Event, error) {
	return s.r.Delete(uid)
}
