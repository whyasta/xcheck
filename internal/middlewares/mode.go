package middlewares

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

func LocalModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer utils.ResponseHandler(c)
		if config.GetAppConfig().APP_ENV != "local" {
			utils.PanicException(response.InvalidRequest, errors.New("Service Unavailable - only local mode").Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
