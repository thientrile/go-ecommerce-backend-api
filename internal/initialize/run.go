package initialize

import (
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
	global.Logger.Info("config log ok!!", zap.String("status", "success"))
	// initialize mysql
	InitMysqlC()
	// initialize redis
	InitRedis()
// init kafka
	InitKafka()
	// initialize router
	r := InitRouter()

	r.Run(fmt.Sprintf(":%v", s.Port)) // Start the server on port 8080
}
