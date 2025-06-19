package main

import (
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/initialize"
)

func main() {
	// Initialize configuration and logger
	initialize.LoadConfig()
	if global.Config.Server.Mode == "dev" {
		initialize.InitLogger()
	}

	// Test the new logger functionality
	global.Logger.PrintStartupBanner(
		"go-ecommerce-backend-api",
		"1.0.0",
		"8002",
		"development",
	)

	// Test initialization steps
	global.Logger.LogInitStart("Database Connection")
	global.Logger.LogDBConnection("MySQL", "localhost:3306", true, nil)
	global.Logger.LogInitStep("MySQL Database", true, nil)

	global.Logger.LogInitStart("Redis Cache")
	global.Logger.LogInitStep("Redis Cache", true, nil)

	global.Logger.LogInitStart("Kafka Message Queue")
	global.Logger.LogInitStep("Kafka Message Queue", true, nil)

	// Test shutdown
	global.Logger.LogShutdown("Test completed")
}
