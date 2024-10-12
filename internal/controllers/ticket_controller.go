package controllers

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type TicketController struct {
	service       *services.TicketService
	importService *services.ImportService
}

func NewTicketController(service *services.TicketService, importService *services.ImportService) *TicketController {
	return &TicketController{
		service:       service,
		importService: importService,
	}
}

func (r TicketController) ImportTicket(c *gin.Context) {
	defer utils.ResponseHandler(c)
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

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

	tempFile := utils.TempFileName("import-ticket", "ticket", ".csv")
	message := "failed"

	fmt.Println("start import ticket => ", len(files))
	for _, file := range files {
		err := c.SaveUploadedFile(file, tempFile)
		if err != nil {
			utils.PanicException(response.InvalidRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		// Open the CSV file
		csvFile, err := os.Open(tempFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		// skip first row
		if _, err := reader.Read(); err != nil {
			log.Fatal("Error reading the first row:", err)
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}

		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}

		uniqueRecords := make(map[string]bool)

		for _, row := range records {
			// Join the row data as a string (you can also customize key based on specific columns)
			key := fmt.Sprintf("%v", row[0])
			// log.Println(key)
			// Check if the row already exists in the map
			if _, exists := uniqueRecords[key]; !exists {
				uniqueRecords[key] = true
			} else {
				message = fmt.Sprintf("Duplicate barcode found: %s. Please review your CSV file", key)
				utils.PanicException(response.InvalidRequest, message)
				return
			}
		}

		for _, row := range records {
			if exists, _ := r.service.Exist(eventID, row[0]); exists {
				message = fmt.Sprintf("Duplicate barcode detected in the database: %s", row[0])
				utils.PanicException(response.InvalidRequest, message)
				return
			}

			_, err := r.service.ValidateRecord(eventID, row)
			if err != nil {
				message = fmt.Sprintf("Error validating record: %s", err.Error())
				utils.PanicException(response.InvalidRequest, message)
				return
			}
		}

		// log.Println(eventID)

		// upload to minio if valid
		bucketName := config.GetAppConfig().MinioBucket
		_, err = r.importService.UploadToMinio(c, bucketName, file, tempFile)
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
			EventID:        &eventID,
			Type:           1, // import ticket
		})
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}

		message = fmt.Sprintf("Uploaded successfully %d files", len(files))
		_, _, err = r.importService.DoImportTicketJob(importFile.ID, eventID, true)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, message, utils.Null()))
}

func (r TicketController) GetImport(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})
	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventID),
	})
	filter = append(filter, utils.Filter{
		Property:  "type",
		Operation: "=",
		Value:     "1",
	})

	rows, count, err := r.service.GetImport(pageParams, filter, sort)

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

func (r TicketController) GetImportDetail(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	importID, err := strconv.Atoi(c.Param("importId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"assign_status"})
	filter = append(filter, utils.Filter{
		Property:  "import_id",
		Operation: "=",
		Value:     strconv.Itoa(importID),
	})
	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventID),
	})
	filter = append(filter, utils.Filter{
		Property:  "assign_status",
		Operation: "=",
		Value:     "1",
	})

	rows, count, err := r.service.GetImportDetail(pageParams, filter, sort)

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

func (r TicketController) GetAll(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})
	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventID),
	})
	filter = append(filter, utils.Filter{
		Property:  "assign_status",
		Operation: "=",
		Value:     "1",
	})

	rows, count, err := r.service.GetFilteredTickets(pageParams, filter, sort)

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

func (r TicketController) Ticket(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var request *dto.TicketRequest

	c.Next()
	c.BindJSON(&request)

	validate := validator.New(validator.WithRequiredStructEnabled())
	en := en.New()
	UniversalTranslator = ut.New(en, en)
	trans, _ := UniversalTranslator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	err = validate.Struct(request)
	if err != nil {
		fmt.Println(err)
		// errors := err.(validator.ValidationErrors)

		// errors := utils.TranslateError(err, trans)
		validatorErrs := err.(validator.ValidationErrors)
		var errors []error
		for _, e := range validatorErrs {
			translatedErr := fmt.Errorf(e.Translate(trans))
			errors = append(errors, translatedErr)
		}
		// fmt.Println(errsEn)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	data, err := r.service.Ticket(int64(eventID), request.OrderID)

	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

/*
func (r TicketController) Import(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var role *models.UserRole

	c.Next()
	c.BindJSON(&role)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(role)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateRole(role)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
*/
