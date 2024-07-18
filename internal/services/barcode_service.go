package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"log"
	"strconv"
)

type BarcodeService struct {
	r repositories.BarcodeRepository
}

func NewBarcodeService(r repositories.BarcodeRepository) *BarcodeService {
	return &BarcodeService{r}
}

func (s *BarcodeService) CreateBarcode(data *models.Barcode) (models.Barcode, error) {
	return s.r.Save(data)
}

func (s *BarcodeService) UpdateBarcode(eventId int64, id int64, data *map[string]interface{}) (models.Barcode, error) {
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
	rows, _, _ := s.r.FindAll([]string{}, utils.NewPaginate(1, 0), filters)

	if len(rows) == 0 {
		return models.Barcode{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *BarcodeService) GetAllBarcodes(pageParams *utils.Paginate, filters []utils.Filter) ([]models.Barcode, int64, error) {
	return s.r.FindAll([]string{"EventAssignment"}, pageParams, filters)
}

func (s *BarcodeService) GetBarcodeByID(uid int64) (models.Barcode, error) {
	return s.r.FindByID(uid)
}

func (s *BarcodeService) Delete(uid int64) (models.Barcode, error) {
	return s.r.Delete(uid)
}

func (s *BarcodeService) CheckBarcode(barcode string) (bool, error) {
	log.Println("start check barcode")
	_, count, _ := s.r.FindAll([]string{"EventAssignment"}, utils.NewPaginate(10, 1), *utils.NewFilters([]utils.Filter{
		{
			Property:  "barcode",
			Operation: "=",
			Value:     barcode,
		},
	}))

	log.Println(count)
	if count <= 0 {
		return false, errors.New("invalid barcode")
	}

	return true, nil
}

func (s *BarcodeService) AssignBarcodes(importId int64, assignId int64) (bool, error) {
	count, err := s.r.AssignBarcodes(importId, assignId)

	if err != nil || count <= 0 {
		return false, err
	}

	return true, nil
}
