package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.POST("", controllers.UserController.CreateUser)
		group.GET("", controllers.UserController.GetAllUser)
		group.GET("/sync", controllers.UserController.GetAllUserSync)
		group.GET("/:id", controllers.UserController.GetUserByID)
		group.POST("/:id", controllers.UserController.UpdateUser)
	}
}
