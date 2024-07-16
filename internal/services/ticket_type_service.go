package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
)

type TicketTypeService struct {
	r repositories.TicketTypeRepository
}

func NewTicketTypeService(r repositories.TicketTypeRepository) *TicketTypeService {
	return &TicketTypeService{r}
}

func (s *TicketTypeService) CreateTicketType(role *models.TicketType) (models.TicketType, error) {
	return s.r.Save(role)
}
func (s *TicketTypeService) GetAllTicketTypes(pageParams *utils.Paginate, filters []utils.Filter) ([]models.TicketType, int64, error) {
	return s.r.FindAll(pageParams, filters)
}

func (s *TicketTypeService) GetTicketTypeByID(uid int64) (models.TicketType, error) {
	return s.r.FindByID(uid)
}

func (s *TicketTypeService) Delete(uid int64) (models.TicketType, error) {
	return s.r.Delete(uid)
}
