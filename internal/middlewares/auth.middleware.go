package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/pkg/utils/auth"
)

type contextKey string

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// GET request url path
		url := c.Request.URL.Path
		log.Println("Request URL:", url)
		// check header Authorization
		jwtToken, ok := auth.ExtractBearerToken(c)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{
				"error":       "Unauthorized",
				"description": "No valid token provided",
			})
			return
		}

		// validate token
		claims, err := auth.ValidateToken(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error":       "Unauthorized",
				"description": "Invalid token",
			})
			return
		}
		// Optionally, set claims in context for downstream handlers
		c.Set("subjectUUID", claims.Subject)		
		c.Set("claims", claims)
		c.Next()
	}
}
