package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Total Errors -> %d", len(c.Errors))

		if len(c.Errors) <= 0 {
			c.Next()
			return
		}

		for _, err := range c.Errors {
			fmt.Printf("Error -> %+v\n", err)
		}
		c.JSON(http.StatusInternalServerError, "")
	}
}
