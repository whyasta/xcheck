package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.POST("", controllers.UserController.CreateUser)
		group.GET("", controllers.UserController.GetAllUser)
		group.GET("/:id", controllers.UserController.GetUserByID)
		group.POST("/:id", controllers.UserController.UpdateUser)
	}
}
