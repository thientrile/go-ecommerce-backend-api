package middlewares

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// ColorLogger returns a gin.LoggerWithFormatter middleware with colored output
func ColorLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = getStatusColor(param.StatusCode)
			methodColor = getMethodColor(param.Method)
			resetColor = "\033[0m"
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}

		// Custom format with emojis and colors
		emoji := getStatusEmoji(param.StatusCode)

		return fmt.Sprintf("%s %s[GIN]%s %v | %s%3d%s | %13v | %15s | %s%-7s%s %s%#v%s\n%s",
			emoji,
			color.CyanString(""),
			resetColor,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			color.WhiteString(""),
			param.Path,
			resetColor,
			param.ErrorMessage,
		)
	})
}

// ColorLoggerToFile returns a gin.LoggerWithFormatter that writes to both console and file
func ColorLoggerToFile() gin.HandlerFunc {
	// File output without colors
	logFile, _ := os.OpenFile("storage/logs/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// Use MultiWriter for dual output
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile)

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			var statusColor, methodColor, resetColor string
			if param.IsOutputColor() {
				statusColor = getStatusColor(param.StatusCode)
				methodColor = getMethodColor(param.Method)
				resetColor = "\033[0m"
			}

			if param.Latency > time.Minute {
				param.Latency = param.Latency - param.Latency%time.Second
			}

			emoji := getStatusEmoji(param.StatusCode)

			return fmt.Sprintf("%s %s[GIN]%s %v | %s%3d%s | %13v | %15s | %s%-7s%s %s%#v%s\n%s",
				emoji,
				color.CyanString(""),
				resetColor,
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				color.WhiteString(""),
				param.Path,
				resetColor,
				param.ErrorMessage,
			)
		},
		Output: gin.DefaultWriter,
	})
}

// getStatusColor returns the appropriate color for HTTP status codes
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // Success - Green
	case code >= 300 && code < 400:
		return "\033[34m" // Redirect - Blue
	case code >= 400 && code < 500:
		return "\033[33m" // Client Error - Yellow
	default:
		return "\033[31m" // Server Error - Red
	}
}

// getMethodColor returns the appropriate color for HTTP methods
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Blue
	case "POST":
		return "\033[32m" // Green
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[36m" // Cyan
	case "HEAD":
		return "\033[35m" // Magenta
	case "OPTIONS":
		return "\033[37m" // White
	default:
		return "\033[0m" // Reset
	}
}

// getStatusEmoji returns emoji based on HTTP status code
func getStatusEmoji(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "✅" // Success
	case code >= 300 && code < 400:
		return "🔀" // Redirect
	case code >= 400 && code < 500:
		return "⚠️ " // Client Error
	case code >= 500:
		return "❌" // Server Error
	default:
		return "ℹ️ " // Info
	}
}

// StartupBanner displays a colorful startup banner
func StartupBanner(port string, domain string) {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════════════════╗")
	color.Cyan("║                                                              ║")
	color.Cyan("║           🚀 GO ECOMMERCE BACKEND API 🚀                     ║")
	color.Cyan("║                                                              ║")
	color.Cyan("╠══════════════════════════════════════════════════════════════╣")
	color.Cyan("║                                                              ║")
	color.Green("║  🌟 Server Status: RUNNING                                   ║")

	// Format port with proper padding
	portLine := fmt.Sprintf("║  🌐 Port: %-51s║", port)
	color.Yellow(portLine)

	if domain != "" {
		// Format domain with proper padding
		domainText := fmt.Sprintf("http://%s", domain)
		domainLine := fmt.Sprintf("║  🌍 Domain: %-49s║", domainText)
		color.Blue(domainLine)

		// Format swagger URL with proper padding
		swaggerText := fmt.Sprintf("http://%s/swagger/index.html", domain)
		swaggerLine := fmt.Sprintf("║  📚 Swagger: %-48s║", swaggerText)
		color.Blue(swaggerLine)
	} else {
		// Format swagger URL for localhost with proper padding
		swaggerText := fmt.Sprintf("http://localhost:%s/swagger/index.html", port)
		swaggerLine := fmt.Sprintf("║  📚 Swagger: %-45s║", swaggerText)
		color.Blue(swaggerLine)
	}

	color.Magenta("║  ⚡ Hot Reload: ENABLED                                      ║")
	color.Cyan("║                                                              ║")
	color.Cyan("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	color.Green("🎯 Ready to handle requests!")
	color.Yellow("🔥 Hot reload is active - Edit files to see live changes")
	fmt.Println()
}
