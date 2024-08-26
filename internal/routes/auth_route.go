package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.POST("/signin", controllers.AuthController.Signin).Use(middlewares.DeviceMiddleware())
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.GET("/me", controllers.AuthController.CurrentUser)
		group.POST("/token", controllers.AuthController.Refresh)
		group.POST("/signout", controllers.AuthController.Signout)
	}
}
