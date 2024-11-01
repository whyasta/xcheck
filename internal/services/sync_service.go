package services

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm/clause"
)

type SyncService struct {
	repoBase       repositories.BaseRepository
	repoEvent      repositories.EventRepository
	repoTicketType repositories.TicketTypeRepository
	repoGate       repositories.GateRepository
	repoSession    repositories.SessionRepository
}

func NewSyncService(
	base repositories.BaseRepository, r repositories.EventRepository, r2 repositories.TicketTypeRepository,
	r3 repositories.GateRepository, r4 repositories.SessionRepository,
) *SyncService {
	return &SyncService{base, r, r2, r3, r4}
}

func (s *SyncService) SyncEvents() (utils.APIResponse[map[string]interface{}], int64, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	req, err := HTTPRequest("GET", config.GetAppConfig().CloudBaseURL+"/events?page=1&limit=999999&filter=[{\"prop\":\"status\",\"opr\":\"=\",\"val\":\"1\"}]", nil)
	if err != nil {
		return utils.APIResponse[map[string]interface{}]{}, 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return utils.APIResponse[map[string]interface{}]{}, 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return utils.APIResponse[map[string]interface{}]{}, 0, errors.New("request failed with status: " + res.Status)
	}

	response := &utils.APIResponse[map[string]interface{}]{
		Data: []models.Event{},
	}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return utils.APIResponse[map[string]interface{}]{}, 0, err
	}

	// fmt.Println(response.Data)

	return *response, 0, nil
}

func (s *SyncService) SyncDownloadEventByID(uid int64) error {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	req, err := HTTPRequest("GET", config.GetAppConfig().CloudBaseURL+"/events/"+strconv.Itoa(int(uid)), nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	response := &utils.APIResponse[map[string]interface{}]{
		Data: models.Event{},
	}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return derr
	}
	// b, _ := json.Marshal(response)
	// fmt.Println(string(b))

	m, _ := response.Data.(map[string]interface{})

	// Event
	eventDto := dto.EventRequest{
		ID:        int64(m["id"].(float64)),
		EventName: m["event_name"].(string),
		StartDate: m["start_date"].(string),
		EndDate:   m["end_date"].(string),
		Status:    int(m["status"].(float64)),
	}
	_, err = s.repoEvent.Save(&eventDto)
	if err != nil {
		return err
	}

	// Ticket Types
	var tiketTypes []models.TicketType
	b, _ := json.Marshal(m["ticket_types"])
	if err := json.Unmarshal(b, &tiketTypes); err != nil {
		return err
	}
	if len(tiketTypes) > 0 {
		_, err = s.repoTicketType.BulkSave(&tiketTypes)
		if err != nil {
			return err
		}
	}

	// Gates
	var gates []dto.GateRequestDto
	b, _ = json.Marshal(m["gates"])
	if err := json.Unmarshal(b, &gates); err != nil {
		return err
	}
	if len(gates) > 0 {
		_, err = s.repoGate.BulkSave(&gates)
		if err != nil {
			return err
		}
	}

	// Sessions
	var sessions []models.Session
	b, _ = json.Marshal(m["sessions"])
	if err := json.Unmarshal(b, &sessions); err != nil {
		return err
	}
	if len(sessions) > 0 {
		_, err = s.repoSession.BulkSave(&sessions)
		if err != nil {
			return err
		}
	}

	// Barcodes
	client = &http.Client{
		Timeout: 5 * time.Minute,
	}
	req, err = HTTPRequest("GET", config.GetAppConfig().CloudBaseURL+"/events/"+strconv.Itoa(int(uid))+"/barcodes?page=1&limit=99999999", nil)
	if err != nil {
		return err
	}
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	//fmt.Println(config.GetAppConfig().CloudBaseURL + "/events/" + strconv.Itoa(int(uid)) + "/barcodes?page=1&limit=99999999")
	response = &utils.APIResponse[map[string]interface{}]{
		Data: models.Event{},
	}
	derr = json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return derr
	}

	var barcodes []models.Barcode
	b, _ = json.Marshal(response.Data)
	if err := json.Unmarshal(b, &barcodes); err != nil {
		return err
	}
	if len(barcodes) > 0 {
		// delete unused barcodes
		s.repoBase.GetDB().Where("current_status = ?", "").Select("Gates", "Sessions").Omit("TicketType").Delete(&barcodes)

		err = s.repoBase.GetDB().Table("barcodes").Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
		}).Omit("LatestScan", "History").Create(&barcodes).Error
		if err != nil {
			return err
		}
	}

	// update event last sync time
	err = s.repoBase.GetDB().Table("events").Where("id = ?", eventDto.ID).Update("last_synced_at", time.Now()).Error
	if err != nil {
		return err
	}

	// s.repoBase.GetDB().Table("barcodes").Create(&response.Data)

	return nil
}

func (s *SyncService) SyncUploadEventByID(uid int64) error {
	var barcodes []models.Barcode
	var barcodeLogs []models.BarcodeLog

	err := s.repoBase.GetDB().
		Table("barcodes").
		Preload("TicketType").
		Preload("Gates").
		Preload("Sessions").
		Where("event_id = ?", uid).
		Find(&barcodes).Error
	if err != nil {
		return err
	}

	s.repoBase.GetDB().Table("barcode_logs").Where("event_id = ?", uid).Find(&barcodeLogs)

	type RequestSyncUpload struct {
		Barcodes []models.Barcode    `json:"barcodes"`
		History  []models.BarcodeLog `json:"history"`
	}

	var reqBody RequestSyncUpload
	reqBody.Barcodes = barcodes
	reqBody.History = barcodeLogs

	body, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// fmt.Println(string(body))
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	req, err := HTTPRequest("POST", config.GetAppConfig().CloudBaseURL+"/barcodes/sync/upload", body)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	response := &utils.APIResponse[map[string]interface{}]{
		Data: map[string]interface{}{},
	}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return derr
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status + " - " + response.Message)
	}
	return nil
}

func backgroundTask(data []byte) {
	config.Logger.Info("Uploading from local to cloud")
}

func (s *SyncService) SyncUsers() error {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	req, err := HTTPRequest("GET", config.GetAppConfig().CloudBaseURL+"/users/sync?page=1&limit=999999", nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	response := &utils.APIResponse[map[string]interface{}]{
		Data: []models.User{},
	}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return derr
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status + " - " + response.Message)
	}

	var users []models.User
	b, _ := json.Marshal(response.Data)
	if err := json.Unmarshal(b, &users); err != nil {
		return err
	}
	if len(users) > 0 {
		err = s.repoBase.GetDB().Table("users").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"email", "role_id", "password"}),
		}).Create(&users).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func HTTPRequest(method string, url string, payload []byte) (*http.Request, error) {
	body := []byte(`{
        "username": "admin",
        "password": "gate@BM2024"
    }`)

	req, err := http.NewRequest("POST", config.GetAppConfig().CloudBaseURL+"/auth/signin", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("HTTP Error:", res.Status)
		return nil, errors.New("request failed with status: " + res.Status)
	}

	token := &models.SignedResponse{}
	derr := json.NewDecoder(res.Body).Decode(token)
	if derr != nil {
		return nil, err
	}

	req, err = http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	return req, nil
}
