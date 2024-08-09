package services

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/dto"
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
	s repositories.GateAllocationRepository
}

func NewBarcodeService(r repositories.BarcodeRepository, s repositories.GateAllocationRepository) *BarcodeService {
	return &BarcodeService{r, s}
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
	rows, _, _ := s.r.FindAll([]string{}, utils.NewPaginate(1, 0), filters, []utils.Sort{})

	if len(rows) == 0 {
		return models.Barcode{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *BarcodeService) DownloadBarcodes(pageParams *utils.Paginate, eventId int64, sessionId int64, gateId int64) ([]models.Barcode, int64, error) {
	gateAllocations, _, err := s.s.FindAll(utils.NewPaginate(999999, 1), *utils.NewFilters([]utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventId)),
		},
		{
			Property:  "session_id",
			Operation: "=",
			Value:     strconv.Itoa(int(sessionId)),
		},
		{
			Property:  "gate_id",
			Operation: "=",
			Value:     strconv.Itoa(int(gateId)),
		},
	}), []utils.Sort{})

	if len(gateAllocations) == 0 || err != nil {
		return []models.Barcode{}, 0, errors.New("barcode not found")
	}

	barcodes, count, err := s.r.FindAll([]string{"GateAllocation"}, pageParams, *utils.NewFilters([]utils.Filter{
		{
			Property:  "gateAllocation_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventId)),
		},
	}), []utils.Sort{})

	return barcodes, count, err
}

func (s *BarcodeService) UploadBarcodeLogs(logs *[]dto.BarcodeUploadLogDto) error {
	barcodeLogs := make([]models.BarcodeLog, 0)
	for _, v := range *logs {
		barcodeLogs = append(barcodeLogs, *v.ToEntity())
	}
	return s.r.CreateBulkLog(&barcodeLogs)
}

func (s *BarcodeService) GetAllBarcodes(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error) {
	return s.r.FindAll([]string{"GateAllocation"}, pageParams, filters, sorts)
}

func (s *BarcodeService) GetBarcodeByID(uid int64) (models.Barcode, error) {
	return s.r.FindByID(uid)
}

func (s *BarcodeService) Delete(uid int64) (models.Barcode, error) {
	return s.r.Delete(uid)
}

func (s *BarcodeService) ScanBarcode(userId int64, eventId int64, gateId int64, barcode string, action constant.BarcodeStatus) (bool, models.Barcode, error) {
	fmt.Printf("START SCAN => BARCODE:%s, EVENT:%d, GATE:%d", barcode, eventId, gateId)
	// _, count, _ := s.r.FindAll([]string{"GateAllocation"}, utils.NewPaginate(10, 1), *utils.NewFilters([]utils.Filter{
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

	if result.GateAllocation.EventID != eventId {
		return false, result, errors.New("wrong event")
	}

	if result.GateAllocation.GateID != gateId {
		return false, result, errors.New("wrong gate")
	}

	fmt.Println("now", time.Now())
	fmt.Println("start", result.GateAllocation.Session.SessionStart)
	fmt.Println("end", result.GateAllocation.Session.SessionEnd)

	if time.Now().After(result.GateAllocation.Session.SessionEnd) {
		// update barcode to expired
		s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagExpired})
		return false, result, errors.New("session ended")
	}

	if !utils.TimeIsBetween(time.Now(), result.GateAllocation.Session.SessionStart, result.GateAllocation.Session.SessionEnd) {
		return false, result, errors.New("not in session time")
	}

	// update barcode to valid
	// s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagUsed})
	firstCheckin, err := s.r.CreateLog(eventId, userId, barcode, result.CurrentStatus, action)
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
