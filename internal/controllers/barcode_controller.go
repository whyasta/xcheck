package controllers

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"fmt"
	"net/http"
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
			ImportedAt:     time.Now().Format(time.RFC3339),
			Status:         string(constant.ImportStatusPending),
			ErrorMessage:   "",
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

	valid, err := r.importService.CheckValid(int64(ba.ImportId), int64(ba.ScheduleID))
	if !valid {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	// process assign barcode to event
	_, err = r.barcodeService.AssignBarcodes(int64(ba.ImportId), int64(ba.ScheduleID))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	message := "success"
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, message, utils.Null()))
}

func (r BarcodeController) ScanBarcode(c *gin.Context) {
	defer utils.ResponseHandler(c)

	userId, _, err := utils.ExtractTokenID(c)

	var scan *dto.ScanBarcode

	c.Next()
	c.BindJSON(&scan)

	firstCheckin, result, err := r.barcodeService.ScanBarcode(userId, scan.EventID, scan.GateID, scan.Barcode)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	message := string(result.CurrentStatus)
	status := response.Success

	if firstCheckin {
		status = response.Checkin
	} else if result.CurrentStatus == constant.BarcodeStatusIn {
		status = response.ReCheckin
	} else {
		status = response.Checkout
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, status, message, utils.Null()))
}
