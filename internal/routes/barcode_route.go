package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func BarcodeRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware())
	{
		group.POST("/upload", controllers.BarcodeController.UploadBarcodes)
		group.POST("/download", controllers.BarcodeController.DownloadBarcodes)
		group.POST("/assign", controllers.BarcodeController.AssignBarcodes)
		group.POST("/scan/:action", controllers.BarcodeController.ScanBarcode)
	}
}
