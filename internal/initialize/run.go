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
	global.Logger.Info("config log ok!!", zap.String("status", "success"))
	// initialize mysql

	// InitMysql()
	// chuyển dùng sqlc và goose
	InitMysqlC()
	// init service interface
	InitServiceInterface()
	// initialize redis
	InitRedis()
	// init kafka
	InitKafka()
	// initialize router
	r := InitRouter()

	// init swagger
	r = InitSwagger(r)
	// r.Run(":8002") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run(fmt.Sprintf(":%v", s.Port))

}
