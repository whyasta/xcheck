package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"fmt"
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

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	data, err := r.service.ReportTrafficVisitor(int64(eventID))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportUniqueVisitor(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
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

	data, err := r.service.ReportUniqueVisitor(int64(eventID), ticketTypeIds, gateIds, sessionIds)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportGateIn(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	data, err := r.service.ReportGateIn(int64(eventID))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportRedemptionSummary(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	ticketTypeIds := make([]int64, 0)
	userIds := make([]int64, 0)
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
		} else if key == "user_id" {
			var params []int64

			if err := json.Unmarshal([]byte(value[0]), &params); err != nil {
				utils.PanicException(response.InvalidRequest, err.Error())
			}

			for _, id := range params {
				userIds = append(userIds, int64(id))
			}
			c.Set("userIds", userIds)
		}
	}

	data, err := r.service.ReportRedemptionSummary(int64(eventID), ticketTypeIds, userIds)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r ReportController) ReportRedemptionLog(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var orderBarcode string
	var orderId string
	for key, value := range c.Request.URL.Query() {
		if key == "order_barcode" {
			orderBarcode = value[0]
		} else if key == "order_id" {
			orderId = value[0]
		}
	}

	fmt.Println("orderBarcode", orderBarcode)
	fmt.Println("orderId", orderId)

	data, err := r.service.ReportRedemptionLog(int64(eventID), orderBarcode, orderId)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}
