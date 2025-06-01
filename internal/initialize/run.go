package initialize

import (
	"fmt"

	"go-ecommerce-backend-api.com/global"
)

func Run() {
	// load configuration
	LoadConfig()
	m := global.Config.Msql
	fmt.Println("Configuration loaded successfully", m.Username, m.Password)
	// initialize logger
	InitLogger()
	// initialize mysql
	InitMysql()
	// initialize redis
	InitRedis()

	// initialize router
	r := InitRouter()

	r.Run(":8002") // Start the server on port 8080
}
