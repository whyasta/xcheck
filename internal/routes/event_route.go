package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func EventRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.GET("", controllers.EventController.GetAllEvents)
		group.POST("", controllers.EventController.CreateEvent)
		group.GET("/:id", controllers.EventController.GetEventByID)
		group.POST("/:id", controllers.EventController.UpdateEvent)

		group.GET("/:id/ticket-types", controllers.TicketTypeController.GetAllTicketTypes)
		group.POST("/:id/ticket-types", controllers.TicketTypeController.CreateTicketType)
		group.POST("/:id/ticket-types/:ticketTypeID", controllers.TicketTypeController.UpdateTicketType)

		group.GET("/:id/gates", controllers.GateController.GetAllGates)
		group.POST("/:id/gates", controllers.GateController.CreateGate)
		group.POST("/:id/gates/:gateID", controllers.GateController.UpdateGate)

		group.GET("/:id/sessions", controllers.SessionController.GetAllSessions)
		group.POST("/:id/sessions", controllers.SessionController.CreateSession)
		group.POST("/:id/sessions/:sessionID", controllers.SessionController.UpdateSession)

		group.GET("/:id/barcodes", controllers.BarcodeController.GetEventBarcodes)
		group.POST("/:id/barcodes/import", controllers.BarcodeController.ImportEventBarcodes)

		group.POST("/:id/barcodes/update-bulk", controllers.BarcodeController.UpdateBulkBarcode)
		group.POST("/:id/barcodes/:barcodeID", controllers.BarcodeController.UpdateBarcode)

		// group.GET("/:id/report", controllers.EventController.ReportEvent)

		// group.GET("/:id/gate-allocations", controllers.GateAllocationController.GetAllGateAllocations)
		// group.POST("/:id/gate-allocations", controllers.GateAllocationController.CreateGateAllocation)
		// group.POST("/:id/gate-allocations/:gateAllocationId", controllers.GateAllocationController.UpdateGateAllocation)

		group.GET("/:id/tickets", controllers.TicketController.GetAll)
		group.POST("/:id/tickets/import", controllers.TicketController.ImportTicket)
		group.GET("/:id/tickets/import", controllers.TicketController.GetImport)
		group.GET("/:id/tickets/import/:importId", controllers.TicketController.GetImportDetail)

		group.POST("/:id/tickets/check", controllers.TicketController.Check)
		group.POST("/:id/tickets/redeem", controllers.TicketController.Redeem)
	}
}
