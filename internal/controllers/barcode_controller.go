package controllers

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	files := form.File["files"]

	if len(files) == 0 {
		utils.PanicException(constant.InvalidRequest, "file not found")
		return
	}

	log.Println("start => ", len(files))
	importFile, err := r.importService.CreateImport(&models.Import{
		FileName:     files[0].Filename,
		ImportedAt:   time.Now().Format(time.RFC3339),
		Status:       string(models.ImportStatusPending),
		ErrorMessage: "",
	})
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	for _, file := range files {
		filename := filepath.Join("files", file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			_, _ = r.importService.UpdateStatusImport(importFile.ID, string(models.ImportStatusFailed), err.Error())
			utils.PanicException(constant.InvalidRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}

	message := fmt.Sprintf("Uploaded successfully %d files", len(files))
	_, err = r.importService.DoImportJob(importFile.ID)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, message, utils.Null()))
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
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	valid, err := r.importService.CheckValid(int64(ba.ImportId), int64(ba.EventAssignmentID))
	if !valid {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	// process assign barcode to event
	_, err = r.barcodeService.AssignBarcodes(int64(ba.ImportId), int64(ba.EventAssignmentID))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	message := "success"
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, message, utils.Null()))
}

func (r BarcodeController) CheckBarcode(c *gin.Context) {
	defer utils.ResponseHandler(c)

	_, err := r.barcodeService.CheckBarcode(c.Param("barcode"))
	if err != nil {
		utils.PanicException(constant.DataNotFound, err.Error())
		return
	}
	message := "success"
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, message, utils.Null()))
}
