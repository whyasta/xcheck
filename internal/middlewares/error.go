package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Errors) <= 0 {
			c.Next()
			return
		}

		fmt.Printf("Total Errors -> %d\n", len(c.Errors))
		for _, err := range c.Errors {
			fmt.Printf("Error -> %+v\n", err)
		}
		c.JSON(http.StatusInternalServerError, "")
	}
}
