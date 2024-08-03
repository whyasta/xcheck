package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func EventRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.GET("", controllers.EventController.GetAllEvents)
		group.POST("", controllers.EventController.CreateEvent)
		group.GET("/:id", controllers.EventController.GetEventByID)
		group.POST("/:id", controllers.EventController.UpdateEvent)

		group.GET("/:id/ticket-types", controllers.TicketTypeController.GetAllTicketTypes)
		group.POST("/:id/ticket-types", controllers.TicketTypeController.CreateTicketType)
		group.POST("/:id/ticket-types/:ticketTypeId", controllers.TicketTypeController.UpdateTicketType)

		group.GET("/:id/gates", controllers.GateController.GetAllGates)
		group.POST("/:id/gates", controllers.GateController.CreateGate)
		group.POST("/:id/gates/:gateId", controllers.GateController.UpdateGate)

		group.GET("/:id/sessions", controllers.SessionController.GetAllSessions)
		group.POST("/:id/sessions", controllers.SessionController.CreateSession)
		group.POST("/:id/sessions/:sessionId", controllers.SessionController.UpdateSession)

		group.GET("/:id/schedules", controllers.ScheduleController.GetAllSchedules)
		group.POST("/:id/schedules", controllers.ScheduleController.CreateSchedule)
		group.POST("/:id/schedules/:scheduleId", controllers.ScheduleController.UpdateSchedule)
	}
}
