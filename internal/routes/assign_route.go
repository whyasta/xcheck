package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AssignRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.GET("", controllers.EventController.GetAllEvents)
		group.POST("", controllers.EventController.CreateEvent)
		group.GET("/:id", controllers.EventController.GetEventByID)
		group.POST("/:id", controllers.EventController.UpdateEvent)
	}
}
