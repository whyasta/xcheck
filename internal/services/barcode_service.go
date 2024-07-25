package services

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"strconv"
	"time"
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
	return s.r.FindAll([]string{"Schedule"}, pageParams, filters)
}

func (s *BarcodeService) GetBarcodeByID(uid int64) (models.Barcode, error) {
	return s.r.FindByID(uid)
}

func (s *BarcodeService) Delete(uid int64) (models.Barcode, error) {
	return s.r.Delete(uid)
}

func (s *BarcodeService) ScanBarcode(userId int64, eventId int64, gateId int64, barcode string) (bool, models.Barcode, error) {
	fmt.Printf("START SCAN => BARCODE:%s, EVENT:%d, GATE:%d", barcode, eventId, gateId)
	// _, count, _ := s.r.FindAll([]string{"Schedule"}, utils.NewPaginate(10, 1), *utils.NewFilters([]utils.Filter{
	// 	{
	// 		Property:  "barcode",
	// 		Operation: "=",
	// 		Value:     barcode,
	// 	},
	// }))

	// fmt.Println(count)
	// if count <= 0 {
	// 	return false, errors.New("invalid barcode")
	// }

	result, err := s.r.Scan(barcode)
	if err != nil {
		return false, result, err
	}

	if result.Schedule.EventID != eventId {
		return false, result, errors.New("wrong event")
	}

	if result.Schedule.GateID != gateId {
		return false, result, errors.New("wrong gate")
	}

	fmt.Println("now", time.Now())
	fmt.Println("start", result.Schedule.Session.SessionStart)
	fmt.Println("end", result.Schedule.Session.SessionEnd)

	if time.Now().After(result.Schedule.Session.SessionEnd) {
		// update barcode to expired
		s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagExpired})
		return false, result, errors.New("session ended")
	}

	if !utils.TimeIsBetween(time.Now(), result.Schedule.Session.SessionStart, result.Schedule.Session.SessionEnd) {
		return false, result, errors.New("not in session time")
	}

	// update barcode to valid
	// s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagUsed})
	firstCheckin, err := s.r.CreateLog(userId, barcode, result.CurrentStatus)
	if err != nil {
		return false, result, err
	}

	result, err = s.r.FindByID(result.ID)

	return firstCheckin, result, err
}

func (s *BarcodeService) AssignBarcodes(importId int64, assignId int64, ticketTypeId int64) (bool, error) {
	count, err := s.r.AssignBarcodes(importId, assignId, ticketTypeId)

	if err != nil || count <= 0 {
		return false, err
	}

	return true, nil
}
