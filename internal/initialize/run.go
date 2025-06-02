package initialize

import (
	"fmt"

	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

func Run() {
	// load configuration
	LoadConfig()
	m := global.Config.Mysql
	fmt.Println("Configuration loaded successfully", m.Username, m.Password)
	// initialize logger
	InitLogger()
	global.Logger.Info("config log ok!!", zap.String("status", "success"))
	// initialize mysql
	InitMysql()
	// initialize redis
	InitRedis()

	// initialize router
	r := InitRouter()

	r.Run(":8002") // Start the server on port 8080
}
