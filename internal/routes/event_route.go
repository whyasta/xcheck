package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func EventRoutes(eventGroup *gin.RouterGroup, controllers *controllers.Controller) {
	eventGroup.Use(middlewares.AuthMiddleware())
	{
		eventGroup.GET("/", controllers.EventController.GetAllEvents)
		eventGroup.POST("/", controllers.EventController.CreateEvent)
		eventGroup.GET("/:id", controllers.EventController.GetEventByID)
		eventGroup.POST("/:id", controllers.EventController.UpdateEvent)

		eventGroup.GET("/:id/ticket-types", controllers.TicketTypeController.GetAllTicketTypes)
		eventGroup.POST("/:id/ticket-types", controllers.TicketTypeController.CreateTicketType)
		eventGroup.POST("/:id/ticket-types/:ticketTypeId", controllers.TicketTypeController.UpdateTicketType)

		eventGroup.GET("/:id/gates", controllers.GateController.GetAllGates)
		eventGroup.POST("/:id/gates", controllers.GateController.CreateGate)
		eventGroup.POST("/:id/gates/:gateId", controllers.GateController.UpdateGate)

		eventGroup.GET("/:id/sessions", controllers.SessionController.GetAllSessions)
		eventGroup.POST("/:id/sessions", controllers.SessionController.CreateSession)
		eventGroup.POST("/:id/sessions/:sessionId", controllers.SessionController.UpdateSession)
	}
}
