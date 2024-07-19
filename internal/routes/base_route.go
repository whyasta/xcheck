package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func BaseRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.GET("/", controllers.HealthController.Init)
	group.GET("/healthcheck", controllers.HealthController.Status)

	group.POST("/signin", controllers.AuthController.Signin)
	group.Use(middlewares.AuthMiddleware())
	{
		// single get
		group.GET("/ticket-types/:ticketTypeId", controllers.TicketTypeController.GetTicketTypeByID)
		group.GET("/gates/:gateId", controllers.GateController.GetGateByID)
		group.GET("/sessions/:sessionId", controllers.SessionController.GetSessionByID)
	}
}
