package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

// LoggerMiddleware tạo middleware để log HTTP requests với format đẹp
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log với format đẹp
		global.Logger.LogHTTPRequest(
			method,
			path,
			clientIP,
			statusCode,
			duration.String(),
		)

		// Log thêm thông tin chi tiết nếu có lỗi
		if statusCode >= 400 {
			if len(c.Errors) > 0 {
				global.Logger.Error("Request completed with errors",
					zap.String("method", method),
					zap.String("path", path),
					zap.String("ip", clientIP),
					zap.Int("status", statusCode),
					zap.Duration("duration", duration),
					zap.Int("bodySize", bodySize),
					zap.String("errors", c.Errors.String()),
				)
			}
		}
	}
}
