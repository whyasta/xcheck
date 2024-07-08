package server

import (
	"bigmind/xcheck-be/controllers"
	"bigmind/xcheck-be/middlewares"
	"bigmind/xcheck-be/services"
	"bigmind/xcheck-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		auth.POST("/signin", controllers.UserController.Signin)
		auth.POST("/signout", middlewares.AuthMiddleware(), controllers.UserController.Signout)
		auth.GET("/me", middlewares.AuthMiddleware(), controllers.UserController.CurrentUser)
	}

	authorized := router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		userGroup := authorized.Group("users")
		{
			userGroup.POST("/", controllers.UserController.CreateUser)
			userGroup.GET("/", controllers.UserController.GetAllUser)
			userGroup.GET("/:id", controllers.UserController.GetUserByID)
		}

		userRoleGroup := authorized.Group("roles")
		{
			userRoleGroup.POST("/", controllers.RoleController.CreateRole)
			userRoleGroup.GET("/", controllers.RoleController.GetAllRole)
			userRoleGroup.GET("/:id", controllers.RoleController.GetRoleByID)
		}

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
