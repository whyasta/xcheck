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
)

type SessionController struct {
	service *services.SessionService
}

func NewSessionController(service *services.SessionService) *SessionController {
	return &SessionController{
		service: service,
	}
}

// swagger:route POST /sessions Session createSession
// Create Session
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r SessionController) CreateSession(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var session *models.Session
	var bulk []models.Session

	jsons := make([]byte, c.Request.ContentLength)
	if _, err := c.Request.Body.Read(jsons); err != nil {
		if err.Error() != "EOF" {
			return
		}
	}

	if err := json.Unmarshal(jsons, &bulk); err != nil {
		bulk = nil
		if err := json.Unmarshal(jsons, &session); err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
	}

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
		result, err := r.service.CreateBulkSession(&bulk)
		if err != nil {
			utils.PanicException(response.InvalidRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
		return
	}

	/*
		c.Next()
		c.BindJSON(&session)*/

	session.EventID = int64(eventId)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(session)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateSession(session)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}

// swagger:route GET /sessions Session getSessionList
// Get Session list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r SessionController) GetAllSessions(c *gin.Context) {
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

	rows, count, err := r.service.GetAllSessions(pageParams, filter, sort)

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

// swagger:route GET /sessions/{id} Session getSession
// Get Session by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r SessionController) GetSessionByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.Session
	user, err = r.service.GetSessionByID(int64(uid))
	if err != nil {
		utils.PanicException(response.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

// DeleteSession swagger:route DELETE /sessions/{id} Session deleteSession
// Delete Session by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (r SessionController) DeleteSession(c *gin.Context) {
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

func (r SessionController) UpdateSession(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	uid, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var session *models.Session
	var request = make(map[string]interface{})

	c.Next()
	c.BindJSON(&request)

	request["event_id"] = int64(eventId)
	err = utils.Decode(request, &session)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	fmt.Println("session", session)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(session)
	if err != nil {
		fmt.Println("err", err)
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	fmt.Println("request", request)

	result, err := r.service.UpdateSession(int64(eventId), int64(uid), &request)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
