package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"strconv"
)

type GateService struct {
	r repositories.GateRepository
}

func NewGateService(r repositories.GateRepository) *GateService {
	return &GateService{r}
}

func (s *GateService) CreateGate(data *models.Gate) (models.Gate, error) {
	return s.r.Save(data)
}

func (s *GateService) UpdateGate(eventId int64, id int64, data *map[string]interface{}) (models.Gate, error) {
	var filters = []utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventId)),
		},
		{
			Property:  "id",
			Operation: "=",
			Value:     strconv.Itoa(int(id)),
		},
	}
	rows, _, _ := s.r.FindAll(utils.NewPaginate(1, 0), filters)

	if len(rows) == 0 {
		return models.Gate{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *GateService) GetAllGates(pageParams *utils.Paginate, filters []utils.Filter) ([]models.Gate, int64, error) {
	return s.r.FindAll(pageParams, filters)
}

func (s *GateService) GetGateByID(uid int64) (models.Gate, error) {
	return s.r.FindByID(uid)
}

func (s *GateService) Delete(uid int64) (models.Gate, error) {
	return s.r.Delete(uid)
}