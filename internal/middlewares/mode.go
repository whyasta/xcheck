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
		if config.GetAppConfig().AppEnv != "local" {
			utils.PanicException(response.InvalidRequest, errors.New("service unavailable - only local mode").Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
