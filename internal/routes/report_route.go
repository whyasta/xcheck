package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.GET("/traffic-visitor", controllers.ReportController.ReportTrafficVisitor)
		group.GET("/unique-visitor", controllers.ReportController.ReportUniqueVisitor)
		group.GET("/gate-in", controllers.ReportController.ReportGateIn)
	}
}
