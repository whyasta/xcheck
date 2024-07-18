package controllers

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"log"
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var session *models.Session

	c.Next()
	c.BindJSON(&session)

	session.EventID = int64(eventId)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(session)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateSession(session)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	pageParams, filter := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})

	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventId),
	})

	rows, count, err := r.service.GetAllSessions(pageParams, filter)

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

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, constant.Success, "", rows, &meta))
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
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user models.Session
	user, err = r.service.GetSessionByID(int64(uid))
	if err != nil {
		utils.PanicException(constant.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
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

func (r SessionController) UpdateSession(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	uid, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var session *models.Session
	var request = make(map[string]interface{})

	c.Next()
	c.BindJSON(&request)

	request["event_id"] = int64(eventId)
	err = utils.Decode(request, &session)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	log.Println("session", session)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(session)
	if err != nil {
		log.Println("err", err)
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	log.Println("request", request)

	result, err := r.service.UpdateSession(int64(eventId), int64(uid), &request)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
}
