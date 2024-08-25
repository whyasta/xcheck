package services

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/constant/response"
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
	g repositories.GateRepository
	s repositories.SessionRepository
}

func NewBarcodeService(r repositories.BarcodeRepository, s repositories.GateRepository, s2 repositories.SessionRepository) *BarcodeService {
	return &BarcodeService{r, s, s2}
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
	// json_contains(sessions, '1') AND
	// json_contains(gates, '2')
	barcodes, count, err := s.r.FindAllWithRelations(pageParams, *utils.NewFilters([]utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventId)),
		},
		{
			Property:  "barcode_sessions.session_id",
			Operation: "has",
			Value:     strconv.Itoa(int(sessionId)),
		},
		{
			Property:  "barcode_gates.gate_id",
			Operation: "has",
			Value:     strconv.Itoa(int(gateId)),
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
	return s.r.FindAllWithRelations(pageParams, filters, sorts)
}

func (s *BarcodeService) GetBarcodeByID(uid int64) (models.Barcode, error) {
	return s.r.FindByID(uid)
}

func (s *BarcodeService) Delete(uid int64) (models.Barcode, error) {
	return s.r.Delete(uid)
}

func (s *BarcodeService) ScanBarcode(userId int64, eventId int64, gateId int64, barcode string, action constant.BarcodeStatus, device string) (bool, models.BarcodeLog, response.ResponseStatus, error) {
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

	result, responseStatus, err := s.r.Scan(barcode)
	if err != nil {
		return false, models.BarcodeLog{}, responseStatus, err
	}
	if result.EventID != eventId {
		return false, models.BarcodeLog{}, response.EC01, errors.New("EC01 - Barcode " + barcode + " not found!")
	}

	// check gate
	validGate := false
	for _, i := range result.Gates {
		if i.ID == gateId {
			validGate = true
			break
		}
	}

	if !validGate {
		gate, _ := s.g.FindByID(gateId)
		return false, models.BarcodeLog{}, response.EC04, fmt.Errorf("EC04 - "+response.EC05.GetResponseMessage(), gate.GateName)
	}

	// check session
	currentSession := int64(0)
	for _, i := range result.Sessions {
		fmt.Println("check session", i.Sessioname)
		fmt.Println("now", time.Now())
		fmt.Println("start", i.SessionStart)
		fmt.Println("end", i.SessionEnd)

		if time.Now().After(i.SessionEnd) {
			// update barcode to expired
			s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagExpired})
			return false, models.BarcodeLog{}, response.EC04, errors.New("EC04 - Barcode " + barcode + " session has ended")
		}

		if !utils.TimeIsBetween(time.Now(), i.SessionStart, i.SessionEnd) {
			return false, models.BarcodeLog{}, response.EC04, errors.New("EC04 - Barcode " + barcode + " is not within the session time")
		} else {
			currentSession = i.ID
			break
		}
	}

	fmt.Println("current session", currentSession)

	if action == constant.BarcodeStatusOut && result.CurrentStatus == constant.BarcodeStatusNull {
		return false, models.BarcodeLog{}, response.EC11, errors.New("EC11 - Barcode " + barcode + " not checked-in yet!")
	}

	if action == constant.BarcodeStatusIn && result.CurrentStatus == constant.BarcodeStatusIn {
		return false, models.BarcodeLog{}, response.EC03, errors.New("EC03 - Barcode " + barcode + " not allowed to re-enter!")
	}

	// update barcode to valid
	// s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagUsed})
	resultLog, firstCheckin, err := s.r.CreateLog(eventId, userId, gateId, result.TicketTypeID, currentSession, barcode, result.CurrentStatus, action, device)
	if err != nil {
		return false, models.BarcodeLog{}, response.UnknownError, err
	}

	// result, err = s.r.FindByID(result.ID)

	return firstCheckin, resultLog, response.UnknownError, err
}

func (s *BarcodeService) AssignBarcodes(importId int64, assignId int64, ticketTypeId int64) (bool, error) {
	count, err := s.r.AssignBarcodes(importId, assignId, ticketTypeId)

	if err != nil || count <= 0 {
		return false, err
	}

	return true, nil
}
