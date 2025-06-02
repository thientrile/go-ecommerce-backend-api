package initialize

import (
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}