package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
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

type GateController struct {
	service *services.GateService
}

func NewGateController(service *services.GateService) *GateController {
	return &GateController{
		service: service,
	}
}

// swagger:route POST /gates Gate createGate
// Create Gate
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateController) CreateGate(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var event *models.Gate

	c.Next()
	c.BindJSON(&event)

	event.EventID = int64(eventId)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(event)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateGate(event)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

// swagger:route GET /gates Gate getGateList
// Get Gate list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateController) GetAllGates(c *gin.Context) {
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

	rows, count, err := r.service.GetAllGates(pageParams, filter, sort)

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

// swagger:route GET /gates/{id} Gate getGate
// Get Gate by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateController) GetGateByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("gateId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.Gate
	user, err = r.service.GetGateByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

// DeleteGate swagger:route DELETE /gates/{id} Gate deleteGate
// Delete Gate by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateController) DeleteGate(c *gin.Context) {
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

func (r GateController) UpdateGate(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	uid, err := strconv.Atoi(c.Param("gateId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var event *models.Gate
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

	result, err := r.service.UpdateGate(int64(eventId), int64(uid), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
