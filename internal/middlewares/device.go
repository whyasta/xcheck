package middlewares

import (
	"bigmind/xcheck-be/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer utils.ResponseHandler(c)
		ua := c.Request.UserAgent()

		isMobile := strings.Contains(strings.ToLower(ua), "dart")
		if isMobile {
			c.Set("device", "mobile")
		} else {
			c.Set("device", "cms")
		}
	}
}
