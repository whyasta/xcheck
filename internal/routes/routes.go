package routes

import (
	"bigmind/xcheck-be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(
	router *gin.Engine,
	controllers *controllers.Controller,
) *gin.Engine {
	BaseRoutes(router.Group("/"), controllers)
	AuthRoutes(router.Group("/auth"), controllers)
	UserRoutes(router.Group("/users"), controllers)
	RoleRoutes(router.Group("/roles"), controllers)
	EventRoutes(router.Group("/events"), controllers)
	BarcodeRoutes(router.Group("/barcodes"), controllers)
	return router
}
