package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"strconv"
)

type ScheduleService struct {
	r repositories.ScheduleRepository
}

func NewScheduleService(r repositories.ScheduleRepository) *ScheduleService {
	return &ScheduleService{r}
}

func (s *ScheduleService) CreateSchedule(data *dto.ScheduleRequest) (models.Schedule, error) {
	return s.r.Save(data.ToEntity())
}

func (s *ScheduleService) UpdateSchedule(eventId int64, id int64, data *map[string]interface{}) (models.Schedule, error) {
	var filters = []utils.Filter{
		{
			Property:  "schedules.event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventId)),
		},
		{
			Property:  "schedules.id",
			Operation: "=",
			Value:     strconv.Itoa(int(id)),
		},
	}
	rows, _, _ := s.r.FindAll(utils.NewPaginate(1, 0), filters)

	if len(rows) == 0 {
		return models.Schedule{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *ScheduleService) GetAllSchedules(pageParams *utils.Paginate, filters []utils.Filter) ([]models.Schedule, int64, error) {
	return s.r.FindAll(pageParams, filters)
}

func (s *ScheduleService) GetScheduleByID(uid int64) (models.Schedule, error) {
	return s.r.FindByID(uid)
}

func (s *ScheduleService) Delete(uid int64) (models.Schedule, error) {
	return s.r.Delete(uid)
}
