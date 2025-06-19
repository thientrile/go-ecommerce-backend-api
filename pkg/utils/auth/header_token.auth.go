package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractBearerToken(c *gin.Context) (string, bool) {
	//Au
	auHeader := c.GetHeader("Authorization") // Return the token without the "Bearer " prefix

	if strings.HasPrefix(auHeader, "Bearer") {
		// Split the header to get the token
		token := strings.TrimSpace(strings.TrimPrefix(auHeader, "Bearer"))
		if token != "" {
			return token, true
		}
	}
	return "", false

}
