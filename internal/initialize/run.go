package initialize

import (
	// "fmt"

	"fmt"

	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

func Run() {
	// load configuration
	LoadConfig()
	s := global.Config.Server

	// initialize logger
	InitLogger()

	// Display startup banner
	global.Logger.PrintStartupBanner("go-ecommerce-backend-api", s.Version, fmt.Sprintf("%v", s.Port), s.Mode)

	// initialize mysql
	global.Logger.LogInitStart("MySQL Database")
	// InitMysql()
	// chuyển dùng sqlc và goose
	InitMysqlC()
	global.Logger.LogInitStep("MySQL Database", true, nil)

	// init service interface
	global.Logger.LogInitStart("Service Interface")
	InitServiceInterface()
	global.Logger.LogInitStep("Service Interface", true, nil)
	// initialize redis
	global.Logger.LogInitStart("Redis")
	InitRedis()
	global.Logger.LogInitStep("Redis", true, nil)

	// init kafka
	global.Logger.LogInitStart("Kafka")
	InitKafka()
	global.Logger.LogInitStep("Kafka", true, nil)
	// initialize rate limiter
	InitLimiter()
	// initialize router
	global.Logger.LogInitStart("Router")
	r := InitRouter()
	global.Logger.LogInitStep("Router", true, nil)

	// init swagger
	global.Logger.LogInitStart("Swagger Documentation")
	r = InitSwagger(r)
	global.Logger.LogInitStep("Swagger Documentation", true, nil)
	// init prometheus
	global.Logger.LogInitStart("Prometheus Monitoring")
	r = InitPrometheus(r)
	global.Logger.LogInitStep("Prometheus Monitoring", true, nil)
	global.Logger.Info(fmt.Sprintf("🌐 Server listening on port %v", s.Port))
	// r.Run(":8002") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := r.Run(fmt.Sprintf(":%v", s.Port)); err != nil {
		global.Logger.Error("❌ Failed to start server", zap.Error(err))
		global.Logger.LogShutdown("Server start failure")
	}
}
