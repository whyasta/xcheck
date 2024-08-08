package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type TicketTypeController struct {
	service *services.TicketTypeService
}

func NewTicketTypeController(service *services.TicketTypeService) *TicketTypeController {
	return &TicketTypeController{
		service: service,
	}
}

// swagger:route POST /ticket-types TicketType createTicketType
// Create TicketType
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r TicketTypeController) CreateTicketType(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var ticketType *models.TicketType
	var bulk []models.TicketType

	jsons := make([]byte, c.Request.ContentLength)
	if _, err := c.Request.Body.Read(jsons); err != nil {
		if err.Error() != "EOF" {
			return
		}
	}

	if err := json.Unmarshal(jsons, &bulk); err != nil {
		bulk = nil
		if err := json.Unmarshal(jsons, &ticketType); err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

	// c.Next()
	// c.BindJSON(&event)

	if bulk != nil {
		for i, item := range bulk {
			bulk[i].EventID = int64(eventId)
			item.EventID = int64(eventId)

			validate := validator.New(validator.WithRequiredStructEnabled())
			validate.RegisterValidation("date", utils.DateValidation)
			err := validate.Struct(&item)
			if err != nil {
				errors := err.(validator.ValidationErrors)
				utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
				return
			}
		}
		result, err := r.service.CreateBulkTicketType(&bulk)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
		return
	}

	ticketType.EventID = int64(eventId)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(ticketType)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateTicketType(ticketType)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

// swagger:route GET /ticket-types TicketType getTicketTypeList
// Get TicketType list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r TicketTypeController) GetAllTicketTypes(c *gin.Context) {
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

	rows, count, err := r.service.GetAllTicketTypes(pageParams, filter, sort)

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

// swagger:route GET /ticket-types/{id} TicketType getTicketType
// Get TicketType by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r TicketTypeController) GetTicketTypeByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("ticketTypeId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.TicketType
	user, err = r.service.GetTicketTypeByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

// DeleteTicketType swagger:route DELETE /ticket-types/{id} TicketType deleteTicketType
// Delete TicketType by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r TicketTypeController) DeleteTicketType(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	_, err = r.service.Delete(int64(uid))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}

func (r TicketTypeController) UpdateTicketType(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	uid, err := strconv.Atoi(c.Param("ticketTypeId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var event *models.TicketType
	var request = make(map[string]interface{})

	// event.EventID = int64(eventId)

	c.Next()
	c.BindJSON(&request)

	request["event_id"] = int64(eventId)
	mapstructure.Decode(request, &event)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(event)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	fmt.Println(request)

	result, err := r.service.UpdateTicketType(int64(eventId), int64(uid), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
