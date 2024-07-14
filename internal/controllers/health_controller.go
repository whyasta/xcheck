package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func (h HealthController) Init(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h HealthController) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
