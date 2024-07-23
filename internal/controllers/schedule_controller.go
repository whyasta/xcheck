package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type ScheduleController struct {
	service *services.ScheduleService
}

func NewScheduleController(service *services.ScheduleService) *ScheduleController {
	return &ScheduleController{
		service: service,
	}
}

// swagger:route POST /schedules Schedule createSchedule
// Create Schedule
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r ScheduleController) CreateSchedule(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var schedule *dto.ScheduleRequest

	c.Next()
	c.BindJSON(&schedule)

	schedule.EventID = int64(uid)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(schedule)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	fmt.Println(schedule)

	result, err := r.service.CreateSchedule(schedule)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

// swagger:route GET /schedules Schedule getScheduleList
// Get Schedule list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r ScheduleController) GetAllSchedules(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})

	filter = append(filter, utils.Filter{
		Property:  "schedules.event_id",
		Operation: "=",
		Value:     strconv.Itoa(uid),
	})

	rows, count, err := r.service.GetAllSchedules(pageParams, filter)

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

// swagger:route GET /schedules/{id} Schedule getSchedule
// Get Schedule by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r ScheduleController) GetScheduleByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("scheduleId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.Schedule
	user, err = r.service.GetScheduleByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

// DeleteSchedule swagger:route DELETE /schedules/{id} Schedule deleteSchedule
// Delete Schedule by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r ScheduleController) DeleteSchedule(c *gin.Context) {
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

func (r ScheduleController) UpdateSchedule(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	scheduleId, err := strconv.Atoi(c.Param("scheduleId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var schedule *models.Schedule
	var request = make(map[string]interface{})

	// schedule.EventID = int64(uid)

	c.Next()
	c.BindJSON(&request)

	fmt.Println(request)

	request["id"] = int64(scheduleId)
	request["event_id"] = int64(uid)
	mapstructure.Decode(request, &schedule)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(schedule)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	fmt.Println(request)

	result, err := r.service.UpdateSchedule(int64(uid), int64(scheduleId), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
