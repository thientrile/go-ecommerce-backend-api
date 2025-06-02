package global

import (
	"go-ecommerce-backend-api.com/pkg/logger"
	"go-ecommerce-backend-api.com/pkg/setting"
	"gorm.io/gorm"
)

var (
	Config setting.Config
	Logger *logger.LoggerZap
	MDB    *gorm.DB
)
