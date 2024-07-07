package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

// @Summary      Healthcheck
// @Description  Healthcheck
// @Tags         common
// @Produce      json
// @Success      200
// @Router       /healthcheck [get]
func (h HealthController) Init(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h HealthController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
