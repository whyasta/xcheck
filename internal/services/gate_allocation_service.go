package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"strconv"
)

type GateAllocationService struct {
	r repositories.GateAllocationRepository
}

func NewGateAllocationService(r repositories.GateAllocationRepository) *GateAllocationService {
	return &GateAllocationService{r}
}

func (s *GateAllocationService) CreateGateAllocation(data *dto.GateAllocationRequest) (models.GateAllocation, error) {
	return s.r.Save(data.ToEntity())
}

func (s *GateAllocationService) CreateBulkGateAllocation(gates *[]models.GateAllocation) ([]models.GateAllocation, error) {
	return s.r.BulkSave(gates)
}

func (s *GateAllocationService) UpdateGateAllocation(eventID int64, id int64, data *map[string]interface{}) (models.GateAllocation, error) {
	var filters = []utils.Filter{
		{
			Property:  "gate_allocations.event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventID)),
		},
		{
			Property:  "gate_allocations.id",
			Operation: "=",
			Value:     strconv.Itoa(int(id)),
		},
	}
	rows, _, _ := s.r.FindAll(utils.NewPaginate(1, 0), filters, []utils.Sort{})

	if len(rows) == 0 {
		return models.GateAllocation{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *GateAllocationService) GetAllGateAllocations(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.GateAllocation, int64, error) {
	return s.r.FindAll(pageParams, filters, sorts)
}

func (s *GateAllocationService) GetGateAllocationByID(uid int64) (models.GateAllocation, error) {
	return s.r.FindByID(uid)
}

func (s *GateAllocationService) Delete(uid int64) (models.GateAllocation, error) {
	return s.r.Delete(uid)
}
