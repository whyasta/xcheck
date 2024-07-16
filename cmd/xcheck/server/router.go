package server

import (
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"net/http"

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

	health := new(controllers.HealthController)
	// user := controllers.NewUserController(services.UserService)

	controllers := controllers.NewController(services)

	router.GET("/", health.Init)
	router.GET("/healthcheck", health.Status)

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

		authorized.POST("/events", controllers.EventController.CreateEvent)
		authorized.GET("/events/:id", controllers.EventController.GetEventByID)
		authorized.POST("/events/:id", controllers.EventController.UpdateEvent)
		authorized.GET("/events", controllers.EventController.GetAllEvents)
		// authorized.DELETE("/events/:id", controllers.EventController.DeleteEvent)

		// authorized.POST("/events/:event_id/gates", controllers.EventController.CreateEvent)
		// authorized.GET("/events/:event_id/gates/:id", controllers.EventController.GetEventByID)
		// authorized.GET("/events/:event_id", controllers.EventController.GetAllEvents)

		authorized.POST("/ticket-types", controllers.TicketTypeController.CreateTicketType)
		authorized.GET("/ticket-types/:id", controllers.TicketTypeController.GetTicketTypeByID)
		authorized.POST("/ticket-types/:id", controllers.TicketTypeController.UpdateTicketType)
		authorized.GET("/ticket-types", controllers.TicketTypeController.GetAllTicketTypes)

	}

	//router.Use(middlewares.AuthMiddleware())

	// v1 := router.Group("v1")
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		user := new(controllers.UserController)
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": 405, "message": "405 method not allowed"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "404 page not found"})
	})

	router.Use(middlewares.ErrorMiddleware())

	return router
}
