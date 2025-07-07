package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/utils/auth"
)

type contextKey string

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// GET request url path
		url := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log auth attempt
		global.Logger.LogAuth("token_validation_attempt", "unknown", ip, userAgent, false, "Checking authorization header")

		// check header Authorization
		jwtToken, ok := auth.ExtractBearerToken(c)
		if !ok {
			// Log auth failure
			global.Logger.LogAuth("token_missing", "unknown", ip, userAgent, false, "No Bearer token provided")
			global.Logger.LogSecurity("unauthorized_access_attempt", "unknown", ip,
				"Attempted to access protected route without token: "+url, "medium")

			c.AbortWithStatusJSON(401, gin.H{
				"error":       "Unauthorized",
				"description": "No valid token provided",
			})
			return
		}

		// validate token
		claims, err := auth.ValidateToken(jwtToken)
		if err != nil {
			// Log auth failure
			global.Logger.LogAuth("token_invalid", "unknown", ip, userAgent, false, "Invalid JWT token: "+err.Error())
			global.Logger.LogSecurity("invalid_token_used", "unknown", ip,
				"Invalid JWT token used for route: "+url, "medium")

			c.AbortWithStatusJSON(401, gin.H{
				"error":       "Unauthorized",
				"description": "Invalid token",
			})
			return
		}

		// Log successful auth
		userID := claims.Subject
		global.Logger.LogAuth("token_validated", userID, ip, userAgent, true, "JWT token validated successfully")

		// Optionally, set claims in context for downstream handlers
		c.Set("subjectUUID", claims.Subject)
		c.Set("user_id", claims.Subject) // For access logging
		c.Set("claims", claims)
		c.Next()
	}
}
