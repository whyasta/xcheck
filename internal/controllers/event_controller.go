package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type EventController struct {
	service *services.EventService
}

func NewEventController(service *services.EventService) *EventController {
	return &EventController{
		service: service,
	}
}

// swagger:route POST /events Event createEvent
// Create Event
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r EventController) CreateEvent(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var event *dto.EventRequest
	var bulkEvent *[]dto.EventRequest

	jsons := make([]byte, c.Request.ContentLength)
	fmt.Println(c.Request.ContentLength)
	if _, err := c.Request.Body.Read(jsons); err != nil {
		if err.Error() != "EOF" {
			return
		}
	}
	fmt.Println(string(jsons))

	if err := json.Unmarshal(jsons, &bulkEvent); err != nil {
		bulkEvent = nil
		if err := json.Unmarshal(jsons, &event); err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

	/*
		c.Next()
		if err := c.BindJSON(&bulkEvent); err != nil {
			bulkEvent = nil
			if err := c.BindJSON(&event); err != nil {
				utils.PanicException(response.InvalidRequest, err.Error())
				return
			}
		}*/

	if bulkEvent != nil {
		for _, event := range *bulkEvent {
			validate := validator.New(validator.WithRequiredStructEnabled())
			validate.RegisterValidation("date", utils.DateValidation)
			err := validate.Struct(event)
			if err != nil {
				errors := err.(validator.ValidationErrors)
				utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
				return
			}

			log.Println(event)
		}
		result, err := r.service.CreateBulkEvent(bulkEvent)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(event)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateEvent(event)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

func (r EventController) UpdateEvent(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var event *dto.EventUpdateDto
	var request = make(map[string]interface{})

	c.Next()
	c.BindJSON(&request)
	mapstructure.Decode(request, &event)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("date", utils.DateValidation)

	err = validate.Struct(event)
	if err != nil {
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", err))
		// errors := err.(validator.ValidationErrors)
		// utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	res, err := r.service.UpdateEvent(int64(uid), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	row := dto.EventResponse{
		ID:          res.ID,
		EventName:   res.EventName,
		Status:      res.Status,
		StartDate:   res.StartDate.Format("2006-01-02"),
		EndDate:     res.EndDate.Format("2006-01-02"),
		TicketTypes: res.TicketTypes,
		Gates:       res.Gates,
		Sessions:    res.Sessions,
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", row))
}

// swagger:route GET /events Event getEventList
// Get Event list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r EventController) GetAllEvents(c *gin.Context) {
	//pageParams, params := MakePaginationQueryParams(c.Request.URL.Query(), []string{"event_id"})
	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})
	rows, count, err := r.service.GetFilteredEvents(pageParams, filter, sort)

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

// swagger:route GET /events/{id} Event getEvent
// Get Event by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r EventController) GetEventByID(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var event dto.EventResponse
	event, err = r.service.GetEventByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", event))
}

// DeleteEvent swagger:route DELETE /events/{id} Event deleteEvent
// Delete Event by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r EventController) DeleteEvent(c *gin.Context) {
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

func (r EventController) ReportEvent(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	data, err := r.service.Report(int64(uid))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))

}
