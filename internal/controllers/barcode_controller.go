package controllers

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BarcodeController struct {
	importService  *services.ImportService
	barcodeService *services.BarcodeService
}

func NewBarcodeController(importService *services.ImportService, barcodeService *services.BarcodeService) *BarcodeController {
	return &BarcodeController{importService, barcodeService}
}

func (r BarcodeController) GetAllUploads(c *gin.Context) {
	defer utils.ResponseHandler(c)
	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"status"})

	rows, count, err := r.importService.GetAllImports(pageParams, filter, sort)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	meta := utils.MetaResponse{
		PagingInfo: utils.PagingInfo{
			Page:  pageParams.GetPage(count),
			Limit: pageParams.GetLimit(count),
			Total: int(count),
		},
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, response.Success, "", rows, &meta))
}

func (r BarcodeController) GetUploadByID(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	model, err := r.importService.GetImportByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", model))
}

func (r BarcodeController) UploadBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	files := form.File["files"]

	if len(files) == 0 {
		utils.PanicException(response.InvalidRequest, "file not found")
		return
	}

	tempFile := utils.TempFileName("files", "barcode", ".csv")
	message := "failed"

	fmt.Println("start => ", len(files))
	for _, file := range files {
		// filename := filepath.Join("files", file.Filename)
		err := c.SaveUploadedFile(file, tempFile)
		if err != nil {
			utils.PanicException(response.InvalidRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		importFile, err := r.importService.CreateImport(&models.Import{
			FileName:       tempFile,
			UploadFileName: files[0].Filename,
			ImportedAt:     time.Now().Format("2006-01-02 15:04:05"),
			Status:         string(constant.ImportStatusPending),
			StatusMessage:  "",
		})
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}

		message = fmt.Sprintf("Uploaded successfully %d files", len(files))
		_, err = r.importService.DoImportJob(importFile.ID)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, message, utils.Null()))
}

func (r BarcodeController) SyncDownloadBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var dto *dto.BarcodeDownloadDto

	c.Next()
	c.BindJSON(&dto)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dto)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	pageParams, _, _ := MakePageFilterQueryParams(c.Request.URL.Query(), []string{""})
	rows, count, err := r.barcodeService.DownloadBarcodes(pageParams, dto.EventID, dto.SessionID, dto.GateID)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	meta := utils.MetaResponse{
		PagingInfo: utils.PagingInfo{
			Page:  pageParams.GetPage(count),
			Limit: pageParams.GetLimit(count),
			Total: int(count),
		},
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, response.Success, "", rows, &meta))
}

func (r BarcodeController) SyncUploadBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var barcodeLogs []dto.BarcodeUploadLogDto

	jsons := make([]byte, c.Request.ContentLength)
	if _, err := c.Request.Body.Read(jsons); err != nil {
		if err.Error() != "EOF" {
			return
		}
	}

	if err := json.Unmarshal(jsons, &barcodeLogs); err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	// c.Next()
	// c.BindJSON(&dto)

	dto := dto.BarcodeUploadDto{
		Data: barcodeLogs,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dto)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	err = r.barcodeService.UploadBarcodeLogs(&dto.Data)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}

func (r BarcodeController) AssignBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)
	var ba *models.BarcodeAssignment

	c.Next()
	c.BindJSON(&ba)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(ba)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	valid, err := r.importService.CheckValid(int64(ba.ImportId), int64(ba.GateAllocationID))
	if !valid {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	// process assign barcode to event
	_, err = r.barcodeService.AssignBarcodes(int64(ba.ImportId), int64(ba.GateAllocationID), int64(ba.TicketTypeID))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	message := "success"
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, message, utils.Null()))
}

func (r BarcodeController) ScanBarcode(c *gin.Context) {
	defer utils.ResponseHandler(c)

	paramAction := c.Param("action")
	if paramAction != "in" && paramAction != "out" {
		utils.PanicException(response.InvalidRequest, "invalid action")
		return
	}

	action := constant.BarcodeStatusIn
	if paramAction == "in" {
		action = constant.BarcodeStatusIn
	} else if paramAction == "out" {
		action = constant.BarcodeStatusOut
	}

	userId, _, _ := utils.ExtractTokenID(c)

	var scan *dto.ScanBarcode

	c.Next()
	c.BindJSON(&scan)

	if scan.Device == "" {
		scan.Device = c.GetHeader("X-Device")
		if scan.Device == "" {
			scan.Device = "cms"
		}

		val, exist := c.Get("device")
		if exist {
			fmt.Println(val)
			scan.Device = val.(string)
		}
	}

	b, _ := json.Marshal(scan)
	fmt.Println(string(b))

	firstCheckin, result, rc, err := r.barcodeService.ScanBarcode(userId, scan.EventID, scan.GateID, scan.Barcode, action, scan.Device)
	if err != nil {
		//newError = errors.New(rc.GetResponseStatus() + " " + err.Error())
		utils.PanicException(rc, err.Error())
		return
	}

	responseLog := dto.BarcodeLogResponseDto{
		ID:            result.ID,
		Barcode:       result.Barcode,
		EventID:       result.EventID,
		GateID:        result.GateID,
		TicketTypeID:  result.TicketTypeID,
		SessionID:     result.SessionID,
		ScannedBy:     result.ScannedBy,
		ScannedAt:     result.ScannedAt,
		Device:        scan.Device,
		Action:        action,
		CurrentStatus: action,
		Flag:          constant.BarcodeFlagUsed,
	}

	var message string
	// _ = string(result.CurrentStatus)

	// message = string(result.CurrentStatus)
	status := response.Success

	if action == constant.BarcodeStatusIn {
		if firstCheckin {
			message = "RE_CHECKIN"
		} else {
			message = "CHECKIN"
		}
	} else {
		message = "CHECKOUT"
	}

	// else if result.CurrentStatus == constant.BarcodeStatusIn {
	// 	status = response.ReCheckin
	// } else {
	// 	status = response.Checkout
	// }

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, status, message, responseLog))
}

func (r BarcodeController) GetEventBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})

	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventId),
	})

	rows, count, err := r.barcodeService.GetAllBarcodes(pageParams, filter, sort)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	meta := utils.MetaResponse{
		PagingInfo: utils.PagingInfo{
			Page:  pageParams.GetPage(count),
			Limit: pageParams.GetLimit(count),
			Total: int(count),
		},
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, response.Success, "", rows, &meta))
}

func (r BarcodeController) ImportEventBarcodes(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	files := form.File["files"]
	ticketTypeId, _ := strconv.ParseInt(c.PostForm("ticket_type_id"), 10, 64)
	sessions := c.PostForm("sessions")
	gates := c.PostForm("gates")

	if len(files) == 0 {
		utils.PanicException(response.InvalidRequest, "file not found")
		return
	}

	tempFile := utils.TempFileName("files", "barcode", ".csv")
	message := "failed"

	fmt.Println("start => ", len(files))

	err = c.SaveUploadedFile(files[0], tempFile)
	if err != nil {
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	importFile, err := r.importService.CreateImport(&models.Import{
		FileName:       tempFile,
		UploadFileName: files[0].Filename,
		ImportedAt:     time.Now().Format("2006-01-02 15:04:05"),
		Status:         string(constant.ImportStatusPending),
		StatusMessage:  "",
	})
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	message = fmt.Sprintf("Uploaded successfully %d files", len(files))
	_, _, err = r.importService.DoImportJobWithAssign(importFile.ID, eventId, ticketTypeId, sessions, gates)

	for {
		val, _ := config.NewRedis().Get().Do("LLEN", "xcheck:jobs:import_barcode")
		if val.(int64) == 1 {
			fmt.Println("Job inprogress, waiting...")
			time.Sleep(1 * time.Second) // Wait before checking again
		} else {
			// check DB
			importRow, _ := r.importService.GetImportByID(importFile.ID)
			if importRow.Status != string(constant.ImportStatusAssigned) {
				fmt.Println("Job status not assigned, waiting...")
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}

	fmt.Println("Job completed")

	// err = processors.ImportBarcodeJob(job)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	uploadResponse, _ := r.importService.GetUploadResult(importFile.ID)

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, message, uploadResponse))
}
