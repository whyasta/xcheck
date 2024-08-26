package routes

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func BarcodeRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		group.GET("/uploads", controllers.BarcodeController.GetAllUploads)
		group.GET("/uploads/:id", controllers.BarcodeController.GetUploadByID)
		group.POST("/uploads", controllers.BarcodeController.UploadBarcodes)

		group.POST("/assign", controllers.BarcodeController.AssignBarcodes)
		group.POST("/scan/:action", controllers.BarcodeController.ScanBarcode)
		group.POST("/sync/download", controllers.BarcodeController.SyncDownloadBarcodes)
		group.POST("/sync/upload", controllers.BarcodeController.SyncUploadBarcodes)
	}
}
