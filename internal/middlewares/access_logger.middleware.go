package middlewares

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// AccessLoggerMiddleware logs all HTTP requests to access log
func AccessLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract user ID from context (if authenticated)
		userID := "anonymous"
		if param.Keys != nil {
			if uid, exists := param.Keys["user_id"]; exists {
				if uidStr, ok := uid.(string); ok {
					userID = uidStr
				}
			}
		}

		// Log using our custom access logger
		global.Logger.LogAccess(
			param.Method,
			param.Path,
			param.ClientIP,
			param.Request.UserAgent(),
			param.StatusCode,
			param.Latency,
			userID,
		)

		// Return empty string since we handle logging ourselves
		return ""
	})
}

// DetailedAccessLoggerMiddleware provides more detailed logging with request/response bodies
func DetailedAccessLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = rw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Extract user ID
		userID := "anonymous"
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}

		// Log with detailed information
		global.Logger.LogAccess(
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Writer.Status(),
			duration,
			userID,
		)

		// Log request/response bodies for debugging (only in dev mode)
		if global.Config.Server.Mode == "dev" && len(requestBody) > 0 {
			global.Logger.Debug("Request Body",
				zap.String("body", string(requestBody)),
			)
		}
	}
}
