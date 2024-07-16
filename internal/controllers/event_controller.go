package controllers

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	var event *models.Event

	c.Next()
	c.BindJSON(&event)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(event)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateEvent(event)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
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
	pageParams, filter := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})
	rows, count, err := r.service.GetFilteredEvents(pageParams, filter)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	meta := utils.MetaResponse{
		Page:  pageParams.GetPage(count),
		Limit: pageParams.GetLimit(count),
		Total: int(count),
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, constant.Success, "", rows, &meta))
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user models.Event
	user, err = r.service.GetEventByID(int64(uid))
	if err != nil {
		utils.PanicException(constant.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	_, err = r.service.Delete(int64(uid))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", utils.Null()))
}
