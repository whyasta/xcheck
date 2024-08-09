package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
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

type GateAllocationController struct {
	service *services.GateAllocationService
}

func NewGateAllocationController(service *services.GateAllocationService) *GateAllocationController {
	return &GateAllocationController{
		service: service,
	}
}

// swagger:route POST /gateAllocations GateAllocation createGateAllocation
// Create GateAllocation
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateAllocationController) CreateGateAllocation(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	/*
		var gateAllocation *dto.GateAllocationRequest

		c.Next()
		c.BindJSON(&gateAllocation)
	*/

	var gateAllocation *models.GateAllocation
	var bulk []models.GateAllocation

	jsons := make([]byte, c.Request.ContentLength)
	if _, err := c.Request.Body.Read(jsons); err != nil {
		if err.Error() != "EOF" {
			return
		}
	}

	if err := json.Unmarshal(jsons, &bulk); err != nil {
		bulk = nil
		if err := json.Unmarshal(jsons, &gateAllocation); err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

	if bulk != nil {
		for i, item := range bulk {
			log.Println("item", item)
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
		result, err := r.service.CreateBulkGateAllocation(&bulk)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
		return
	}

	gateAllocation.EventID = int64(eventId)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(gateAllocation)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateGateAllocation(&dto.GateAllocationRequest{
		EventID:   gateAllocation.EventID,
		SessionID: gateAllocation.SessionID,
		GateID:    gateAllocation.GateID,
	})
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

// swagger:route GET /gateAllocations GateAllocation getGateAllocationList
// Get GateAllocation list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateAllocationController) GetAllGateAllocations(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})

	filter = append(filter, utils.Filter{
		Property:  "gate_allocations.event_id",
		Operation: "=",
		Value:     strconv.Itoa(uid),
	})

	rows, count, err := r.service.GetAllGateAllocations(pageParams, filter, sort)

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

// swagger:route GET /gateAllocations/{id} GateAllocation getGateAllocation
// Get GateAllocation by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateAllocationController) GetGateAllocationByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("gateAllocationId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.GateAllocation
	user, err = r.service.GetGateAllocationByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

// DeleteGateAllocation swagger:route DELETE /gateAllocations/{id} GateAllocation deleteGateAllocation
// Delete GateAllocation by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r GateAllocationController) DeleteGateAllocation(c *gin.Context) {
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

func (r GateAllocationController) UpdateGateAllocation(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	gateAllocationId, err := strconv.Atoi(c.Param("gateAllocationId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var gateAllocation *models.GateAllocation
	var request = make(map[string]interface{})

	// gateAllocation.EventID = int64(uid)

	c.Next()
	c.BindJSON(&request)

	fmt.Println(request)

	request["id"] = int64(gateAllocationId)
	request["event_id"] = int64(uid)
	mapstructure.Decode(request, &gateAllocation)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(gateAllocation)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	fmt.Println(request)

	result, err := r.service.UpdateGateAllocation(int64(uid), int64(gateAllocationId), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
