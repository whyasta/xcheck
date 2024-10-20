package routes

import (
	"bigmind/xcheck-be/checks"
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/internal/middlewares"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func BaseRoutes(group *gin.RouterGroup, controllers *controllers.Controller) {
	// mysql check
	db, _ := sql.Open("mysql", config.GetDsn())
	mysqlCheck := checks.SQLCheck{SQL: db}

	// redis check
	redisCheck := checks.RedisCheck{Pool: config.NewRedis()}

	failureNotification := checks.FailureNotification{
		Threshold: 3,
		Chan:      make(chan error, 1),
	}

	group.GET("/", controllers.HealthController.Init)
	group.GET("/healthcheck", controllers.HealthController.Status([]checks.Check{mysqlCheck, redisCheck}, &failureNotification))
	group.GET("/timeout", controllers.HealthController.Timeout)

	// group.POST("/signin", controllers.AuthController.Signin)
	group.Use(middlewares.AuthMiddleware(), middlewares.DeviceMiddleware())
	{
		// single get
		group.GET("/ticket-types/:ticketTypeID", controllers.TicketTypeController.GetTicketTypeByID)
		group.GET("/gates/:gateID", controllers.GateController.GetGateByID)
		group.GET("/sessions/:sessionID", controllers.SessionController.GetSessionByID)
		// group.GET("/gate-allocations/:gateAllocationId", controllers.GateAllocationController.GetGateAllocationByID)
	}
}
