package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SyncRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware())
	{
		group.GET("/events", controllers.SyncController.SyncEvents)
		group.POST("/events/:id", controllers.SyncController.SyncEventByID)
	}
}
