package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func BarcodeRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(controllers))
	{
		group.POST("/upload", controllers.BarcodeController.UploadBarcodes)
		group.POST("/assign", controllers.BarcodeController.AssignBarcodes)
		group.POST("/scan", controllers.BarcodeController.ScanBarcode)
	}
}
