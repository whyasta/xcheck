package services

import (
	"bigmind/xcheck-be/internal/dto"
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

func (s *GateService) CreateGate(data *dto.GateRequestDto) (models.Gate, error) {
	return s.r.Save(data)
}

func (s *GateService) CreateBulkGate(gates *[]dto.GateRequestDto) ([]models.Gate, error) {
	return s.r.BulkSave(gates)
}

func (s *GateService) UpdateGate(eventId int64, id int64, data *dto.GateRequestDto) (models.Gate, error) {
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
	rows, _, _ := s.r.FindAll(utils.NewPaginate(1, 0), filters, []utils.Sort{})

	if len(rows) == 0 {
		return models.Gate{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *GateService) GetAllGates(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Gate, int64, error) {
	return s.r.FindAll(pageParams, filters, sorts)
}

func (s *GateService) GetGateByID(uid int64) (models.Gate, error) {
	return s.r.FindByID(uid)
}

func (s *GateService) Delete(uid int64) (models.Gate, error) {
	return s.r.Delete(uid)
}
