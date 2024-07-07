package middlewares

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/token"
	"bigmind/xcheck-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			//c.String(http.StatusUnauthorized, "Unauthorized")
			c.JSON(http.StatusUnauthorized, utils.BuildResponse(http.StatusUnauthorized, constant.Unauthorized, err.Error(), utils.Null()))
			c.Abort()
			return
		}
		c.Next()
	}
}
