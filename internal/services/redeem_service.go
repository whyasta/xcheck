package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
)

type RedeemService struct {
	r repositories.RedeemRepository
}

func NewRedeemService(r repositories.RedeemRepository) *RedeemService {
	return &RedeemService{r}
}

func (s *RedeemService) GetFilteredRedeems(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Redeem, int64, error) {
	result, count, err := s.r.GetFiltered(pageParams, filters, sorts)
	return result, count, err
}

func (s *RedeemService) Redeem(eventID int64, orderID string) (models.Redeem, error) {
	result, err := s.r.FindByOrderID(eventID, orderID)
	return result, err
}
