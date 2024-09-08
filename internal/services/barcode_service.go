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

func (s *BarcodeService) UpdateBarcode(eventID int64, id int64, data *map[string]interface{}) (models.Barcode, error) {
	var filters = []utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventID)),
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

func (s *BarcodeService) DownloadBarcodes(pageParams *utils.Paginate, eventID int64, sessionID int64, gateID int64) ([]models.Barcode, int64, error) {
	// json_contains(sessions, '1') AND
	// json_contains(gates, '2')
	barcodes, count, err := s.r.FindAllWithRelations(pageParams, *utils.NewFilters([]utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventID)),
		},
		{
			Property:  "barcode_sessions.session_id",
			Operation: "has",
			Value:     strconv.Itoa(int(sessionID)),
		},
		{
			Property:  "barcode_gates.gate_id",
			Operation: "has",
			Value:     strconv.Itoa(int(gateID)),
		},
	}), []utils.Sort{})

	return barcodes, count, err
}

func (s *BarcodeService) UploadBarcode(data *[]dto.BarcodeUploadDataDto) error {
	barcodes := make([]models.Barcode, 0)
	for _, v := range *data {
		barcodes = append(barcodes, *v.ToEntity())
	}
	return s.r.CreateBulk(&barcodes)
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

func (s *BarcodeService) ScanBarcode(userID int64, eventID int64, gateID int64, sessionID *int64, barcode string, action constant.BarcodeStatus, device string) (bool, models.BarcodeLog, response.ResponseStatus, error) {
	fmt.Printf("START SCAN => BARCODE:%s, EVENT:%d, GATE:%d", barcode, eventID, gateID)
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

	result, responseStatus, err := s.r.Scan(eventID, barcode)
	if err != nil {
		return false, models.BarcodeLog{}, responseStatus, err
	}
	if result.EventID != eventID {
		return false, models.BarcodeLog{}, response.EC01, errors.New("EC01 - Barcode " + barcode + " not found!")
	}

	// check gate
	validGate := false
	for _, i := range result.Gates {
		if i.ID == gateID {
			validGate = true
			break
		}
	}

	if !validGate {
		gate, _ := s.g.FindByID(gateID)
		return false, models.BarcodeLog{}, response.EC04, fmt.Errorf("EC04 - "+response.EC05.GetResponseMessage(), gate.GateName)
	}

	// check session
	validSession := false
	currentSession := int64(0)

	if sessionID != nil {
		fmt.Println("session param", int64(*sessionID))
		for _, i := range result.Sessions {
			if i.ID == *sessionID {
				fmt.Println("check session", i.Sessioname)
				fmt.Println("now", time.Now())
				fmt.Println("start", i.SessionStart)
				fmt.Println("end", i.SessionEnd)

				if time.Now().After(i.SessionEnd) {
					validSession = false
				}

				if !utils.TimeIsBetween(time.Now(), i.SessionStart, i.SessionEnd) {
					validSession = false
				} else {
					validSession = true
					currentSession = i.ID
					break
				}
				break
			}
			fmt.Println("session available", int64(i.ID))
		}
		if validSession {
			currentSession = int64(*sessionID)
		}
	} else {
		validSession = false
		return false, models.BarcodeLog{}, response.EC04, fmt.Errorf("EC99 - param session id is required")
		/*
			for _, i := range result.Sessions {
				fmt.Println("check session", i.Sessioname)
				fmt.Println("now", time.Now())
				fmt.Println("start", i.SessionStart)
				fmt.Println("end", i.SessionEnd)

				if time.Now().After(i.SessionEnd) {
					// update barcode to expired
					// s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagExpired})
					//return false, models.BarcodeLog{}, response.EC04, errors.New("EC04 - Barcode " + barcode + " session has ended")
					validSession = false
				}

				if !utils.TimeIsBetween(time.Now(), i.SessionStart, i.SessionEnd) {
					validSession = false
					//return false, models.BarcodeLog{}, response.EC04, errors.New("EC04 - Barcode " + barcode + " is not within the session time")
				} else {
					validSession = true
					currentSession = i.ID
					break
				}
			}
		*/
	}

	fmt.Println("validSession", validSession)

	if !validSession {
		return false, models.BarcodeLog{}, response.EC04, fmt.Errorf("EC04 - Barcode " + barcode + " is not within the session time")
	}

	fmt.Println("current session", currentSession)

	if action == constant.BarcodeStatusOut && result.CurrentStatus == constant.BarcodeStatusOut {
		return false, models.BarcodeLog{}, response.EC11, errors.New("EC11 - Barcode " + barcode + " must be checked in first!")
	}

	if action == constant.BarcodeStatusOut && result.CurrentStatus == constant.BarcodeStatusNull {
		return false, models.BarcodeLog{}, response.EC11, errors.New("EC11 - Barcode " + barcode + " not checked-in yet!")
	}

	if action == constant.BarcodeStatusIn && result.CurrentStatus == constant.BarcodeStatusIn {
		return false, models.BarcodeLog{}, response.EC03, errors.New("EC03 - Barcode " + barcode + " not allowed to re-enter!")
	}

	// update barcode to valid
	// s.r.Update(result.ID, &map[string]interface{}{"flag": constant.BarcodeFlagUsed})
	resultLog, firstCheckin, err := s.r.CreateLog(eventID, userID, gateID, result.TicketTypeID, currentSession, barcode, result.CurrentStatus, action, device)
	if err != nil {
		return false, models.BarcodeLog{}, response.UnknownError, err
	}

	// result, err = s.r.FindByID(result.ID)

	return firstCheckin, resultLog, response.UnknownError, err
}

func (s *BarcodeService) AssignBarcodes(importID int64, assignID int64, ticketTypeID int64) (bool, error) {
	count, err := s.r.AssignBarcodes(importID, assignID, ticketTypeID)

	if err != nil || count <= 0 {
		return false, err
	}

	return true, nil
}
