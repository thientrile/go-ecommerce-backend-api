package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "valid-token" {
			response.ErrorResponse(c, response.ErrCodeTokenInvalid, "")
			c.Abort() // Stop the request from proceeding further
			return
		}
		c.Next()
	}
}
