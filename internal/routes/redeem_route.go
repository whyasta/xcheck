package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RedeemRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.GET("", controllers.RedeemController.GetAll)
		group.POST("/redeem", controllers.RedeemController.Redeem)
		// group.POST("/import", controllers.RedeemController.Import)
	}
}
