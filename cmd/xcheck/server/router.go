package server

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"
	"bigmind/xcheck-be/internal/routes"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func NewRouter(services *services.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.Group("/swagger").Handler(ginSwagger.WrapHandler(
	// 	httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("LISTENADDR"))),
	// 	httpSwagger.DeepLinking(true),
	// 	httpSwagger.DocExpansion("none"),
	// 	httpSwagger.DomID("swagger-ui"),
	// )).Methods(http.MethodGet)

	// router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/swagger.yml", "./docs/swagger.yml")
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.GET("/docs", gin.WrapH(sh))

	router.Use(utils.WriterHandler)
	// router.Use(utils.ResponseLogger())
	// router.Use(utils.ResponseHandler())

	controllers := controllers.NewController(services)
	routes.InitRoutes(router, controllers)

	/*
		auth := router.Group("/auth")
		{
			auth.POST("/signin", controllers.AuthController.Signin)
			auth.Use(middlewares.AuthMiddleware())
			{
				auth.GET("/me", controllers.AuthController.CurrentUser)
				auth.POST("/token", controllers.AuthController.Refresh)
				auth.POST("/signout", controllers.AuthController.Signout)
			}
		}

		authorized := router.Group("/")
		authorized.Use(middlewares.AuthMiddleware())
		{
			userGroup := authorized.Group("users")
			{
				userGroup.POST("/", controllers.UserController.CreateUser)
				userGroup.GET("/", controllers.UserController.GetAllUser)
				userGroup.GET("/:id", controllers.UserController.GetUserByID)
				userGroup.POST("/:id", controllers.UserController.UpdateUser)
			}

			userRoleGroup := authorized.Group("roles")
			{
				userRoleGroup.POST("/", controllers.RoleController.CreateRole)
				userRoleGroup.GET("/", controllers.RoleController.GetAllRole)
				userRoleGroup.GET("/:id", controllers.RoleController.GetRoleByID)
				userRoleGroup.POST("/:id", controllers.EventController.UpdateEvent)
			}

			eventGroup := authorized.Group("events")
			{
				eventGroup.GET("/", controllers.EventController.GetAllEvents)
				eventGroup.POST("/", controllers.EventController.CreateEvent)
				eventGroup.GET("/:id", controllers.EventController.GetEventByID)
				eventGroup.POST("/:id", controllers.EventController.UpdateEvent)

				eventGroup.GET("/:id/ticket-types", controllers.TicketTypeController.GetAllTicketTypes)
				eventGroup.POST("/:id/ticket-types", controllers.TicketTypeController.CreateTicketType)
				eventGroup.POST("/:id/ticket-types/:ticketTypeId", controllers.TicketTypeController.UpdateTicketType)

				eventGroup.GET("/:id/gates", controllers.GateController.GetAllGates)
				eventGroup.POST("/:id/gates", controllers.GateController.CreateGate)
				eventGroup.POST("/:id/gates/:gateId", controllers.GateController.UpdateGate)

				eventGroup.GET("/:id/sessions", controllers.SessionController.GetAllSessions)
				eventGroup.POST("/:id/sessions", controllers.SessionController.CreateSession)
				eventGroup.POST("/:id/sessions/:sessionId", controllers.SessionController.UpdateSession)

			}

			// single get
			authorized.GET("/ticket-types/:ticketTypeId", controllers.TicketTypeController.GetTicketTypeByID)
			authorized.GET("/gates/:gateId", controllers.GateController.GetGateByID)
			authorized.GET("/sessions/:sessionId", controllers.SessionController.GetSessionByID)

			// barcodes
			barcodeGroup := authorized.Group("barcodes")
			{
				barcodeGroup.POST("/upload", controllers.BarcodeController.UploadBarcodes)
				barcodeGroup.POST("/assign", controllers.BarcodeController.AssignBarcodes)
				barcodeGroup.POST("/scan", controllers.BarcodeController.ScanBarcode)
			}
		}*/

	//router.Use(middlewares.AuthMiddleware())

	// v1 := router.Group("v1")
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		user := new(controllers.UserController)
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }

	// router.NoMethod(func(c *gin.Context) {
	// 	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": 405, "message": "405 method not allowed"})
	// })

	// router.NoRoute(func(c *gin.Context) {
	// 	c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "404 page not found"})
	// })

	router.Use(middlewares.ErrorMiddleware())

	return router
}
