package examples

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/response"
)

// Example usage trong User Login Controller
func LoginExample(ctx *gin.Context) {
	userID := "user123"
	ip := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	// Log auth attempt
	global.Logger.LogAuth("login_attempt", userID, ip, userAgent, false, "Invalid password")

	// Simulate login logic...
	success := false // assume login failed

	if success {
		// Log successful login
		global.Logger.LogAuth("login_success", userID, ip, userAgent, true, "User logged in successfully")

		// Log business event
		global.Logger.LogBusiness("user_login", userID, "User logged in from web", map[string]interface{}{
			"platform": "web",
			"location": "Vietnam",
		})
	} else {
		// Log failed login
		global.Logger.LogAuth("login_failed", userID, ip, userAgent, false, "Invalid credentials")

		// Log security event for multiple failed attempts
		global.Logger.LogSecurity("multiple_login_failures", userID, ip, "User has 3 failed login attempts", "medium")
	}

	response.ErrorResponse(ctx, response.ErrCodeParamInvalid, "Login failed")
}

// Example usage trong Database operations
func DatabaseOperationExample() {
	start := time.Now()

	// Simulate database operation
	query := "SELECT * FROM users WHERE id = ?"
	operation := "SELECT"
	table := "users"

	// Simulate query execution
	time.Sleep(50 * time.Millisecond)
	duration := time.Since(start)

	// Log successful database operation
	global.Logger.LogDatabase(operation, table, query, duration, 1, nil)

	// Example of failed operation
	// global.Logger.LogDatabase(operation, table, query, duration, 0, errors.New("connection timeout"))
}

// Example usage trong Payment processing
func PaymentExample() {
	userID := "user123"
	orderID := "order456"
	amount := "100.00"
	currency := "USD"
	gateway := "stripe"

	// Log payment attempt
	global.Logger.LogPayment("payment_attempt", userID, orderID, amount, currency, gateway, false, "Processing payment")

	// Simulate payment processing...
	success := true

	if success {
		global.Logger.LogPayment("payment_success", userID, orderID, amount, currency, gateway, true, "Payment processed successfully")

		// Log business event
		global.Logger.LogBusiness("order_paid", userID, "Order payment completed", map[string]interface{}{
			"order_id": orderID,
			"amount":   amount,
			"currency": currency,
			"gateway":  gateway,
		})
	} else {
		global.Logger.LogPayment("payment_failed", userID, orderID, amount, currency, gateway, false, "Insufficient funds")

		// Log security event for suspicious payment
		global.Logger.LogSecurity("suspicious_payment", userID, "", "Multiple failed high-value payments", "high")
	}
}

// Example usage trong System events
func SystemEventExample() {
	// Log system startup
	global.Logger.LogSystem("database", "connection_established", "Connected to MySQL database", "info")

	// Log system error
	global.Logger.LogSystem("redis", "connection_failed", "Failed to connect to Redis server", "error")

	// Log system warning
	global.Logger.LogSystem("memory", "high_usage", "Memory usage above 80%", "warn")
}
