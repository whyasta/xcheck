package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.POST("/", controllers.RoleController.CreateRole)
		group.GET("/", controllers.RoleController.GetAllRole)
		group.GET("/:id", controllers.RoleController.GetRoleByID)
		group.POST("/:id", controllers.RoleController.CreateRole)
	}
}
