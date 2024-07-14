package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Init(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Status describes the Status function in the HealthController.
// @Summary      Healthcheck
// @Description  Healthcheck
// @Tags         common
// @Produce      json
// @Success      200
// @Router       /healthcheck [get]
func (h HealthController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
