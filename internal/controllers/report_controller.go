package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service *services.ReportService
}

func NewReportController(service *services.ReportService) *ReportController {
	return &ReportController{
		service: service,
	}
}

func (r ReportController) ReportTrafficVisitor(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	data, err := r.service.ReportTrafficVisitor(int64(eventId))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportUniqueVisitor(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	ticketTypeIds := make([]int64, 0)
	for key, value := range c.Request.URL.Query() {
		if key == "ticket_type_id" {
			var params []int64

			if err := json.Unmarshal([]byte(value[0]), &params); err != nil {
				utils.PanicException(response.InvalidRequest, err.Error())
			}

			for _, id := range params {
				ticketTypeIds = append(ticketTypeIds, int64(id))
			}
			c.Set("ticketTypeIds", ticketTypeIds)
		}
	}

	gateIds := make([]int64, 0)
	for key, value := range c.Request.URL.Query() {
		if key == "gate_id" {
			var params []int64

			if err := json.Unmarshal([]byte(value[0]), &params); err != nil {
				utils.PanicException(response.InvalidRequest, err.Error())
			}

			for _, id := range params {
				gateIds = append(gateIds, int64(id))
			}
			c.Set("gateIds", gateIds)
		}
	}

	sessionIds := make([]int64, 0)
	for key, value := range c.Request.URL.Query() {
		if key == "session_id" {
			var params []int64

			if err := json.Unmarshal([]byte(value[0]), &params); err != nil {
				utils.PanicException(response.InvalidRequest, err.Error())
			}

			for _, id := range params {
				sessionIds = append(sessionIds, int64(id))
			}
			c.Set("sessionIds", sessionIds)
		}
	}

	data, err := r.service.ReportUniqueVisitor(int64(eventId), ticketTypeIds, gateIds, sessionIds)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportGateIn(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	data, err := r.service.ReportGateIn(int64(eventId))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}
