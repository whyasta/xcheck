package middlewares

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/controllers"
	"bigmind/xcheck-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(controllers *controllers.Controller) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.TokenValid(c)
		if err != nil {
			//c.String(http.StatusUnauthorized, "Unauthorized")
			c.JSON(http.StatusUnauthorized, utils.BuildResponse(http.StatusUnauthorized, response.Unauthorized, err.Error(), utils.Null()))
			c.Abort()
			return
		}
		// if !controllers.AuthController.CheckAuthID(c) {
		// 	fmt.Println("Unauthorized")
		// 	c.JSON(http.StatusUnauthorized, utils.BuildResponse(http.StatusUnauthorized, response.Unauthorized, "Unauthorized", utils.Null()))
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
