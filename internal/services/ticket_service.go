package services

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
)

type TicketService struct {
	r repositories.TicketRepository
	i repositories.ImportRepository
}

func NewTicketService(r repositories.TicketRepository, i repositories.ImportRepository) *TicketService {
	return &TicketService{r, i}
}

func (s *TicketService) GetImport(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error) {
	result, count, err := s.r.GetImport(pageParams, filters, sorts)
	return result, count, err
}

func (s *TicketService) GetImportDetail(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error) {
	result, count, err := s.r.GetImportDetail(pageParams, filters, sorts)
	return result, count, err
}

func (s *TicketService) GetFilteredTickets(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error) {
	result, count, err := s.r.GetFiltered(pageParams, filters, sorts)
	return result, count, err
}

func (s *TicketService) Exist(eventID int64, orderBarcode string) (bool, error) {
	result, err := s.r.Exist(eventID, orderBarcode)
	return result, err
}

func (s *TicketService) ValidateRecord(eventID int64, row []string) (bool, error) {
	result, err := s.r.ValidateRecord(eventID, row)
	return result, err
}

func (s *TicketService) Check(eventID int64, orderBarcode string) (models.Ticket, error) {
	result, err := s.r.FindByBarcode(eventID, orderBarcode)
	return result, err
}

func (s *TicketService) Redeem(eventID int64, data []dto.TicketRedeemDataRequest) ([]models.Ticket, error) {
	result, err := s.r.Redeem(eventID, data)
	return result, err
}
