package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.POST("/signin", controllers.AuthController.Signin)
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.GET("/me", controllers.AuthController.CurrentUser)
		group.POST("/token", controllers.AuthController.Refresh)
		group.POST("/signout", controllers.AuthController.Signout)
	}
}
