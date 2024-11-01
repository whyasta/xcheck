package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SyncController struct {
	service *services.SyncService
}

func NewSyncController(service *services.SyncService) *SyncController {
	return &SyncController{service}
}

func (s SyncController) SyncDownload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Sync download"})
}

func (s SyncController) SyncEvents(c *gin.Context) {
	defer utils.ResponseHandler(c)

	data, _, _ := s.service.SyncEvents()
	c.JSON(http.StatusOK, data)
}

func (s SyncController) SyncDownloadEventByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	err = s.service.SyncDownloadEventByID(int64(uid))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}

func (s SyncController) SyncUploadEventByID(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	err = s.service.SyncUploadEventByID(int64(uid))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}

func (s SyncController) SyncUsers(c *gin.Context) {
	defer utils.ResponseHandler(c)

	err := s.service.SyncUsers()
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}
