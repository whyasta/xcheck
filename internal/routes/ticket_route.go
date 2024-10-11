package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func TicketRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		// group.GET("", controllers.TicketController.GetAll)
	}
}
