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
	mysqlCheck := checks.SqlCheck{Sql: db}

	// redis check
	redisCheck := checks.RedisCheck{Pool: config.NewRedis()}

	failureNotification := checks.FailureNotification{
		Threshold: 3,
		Chan:      make(chan error, 1),
	}

	group.GET("/", controllers.HealthController.Init)
	group.GET("/healthcheck", controllers.HealthController.Status([]checks.Check{mysqlCheck, redisCheck}, &failureNotification))

	// group.POST("/signin", controllers.AuthController.Signin)
	group.Use(middlewares.AuthMiddleware())
	{
		// single get
		group.GET("/ticket-types/:ticketTypeId", controllers.TicketTypeController.GetTicketTypeByID)
		group.GET("/gates/:gateId", controllers.GateController.GetGateByID)
		group.GET("/sessions/:sessionId", controllers.SessionController.GetSessionByID)
		group.GET("/gate-allocations/:gateAllocationId", controllers.GateAllocationController.GetGateAllocationByID)
	}
}
