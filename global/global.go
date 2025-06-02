package global

import (
	"github.com/redis/go-redis/v9"
	"go-ecommerce-backend-api.com/pkg/logger"
	"go-ecommerce-backend-api.com/pkg/setting"
	"gorm.io/gorm"
)

var (
	Config setting.Config
	Logger *logger.LoggerZap
	RDB *redis.Client
	MDB    *gorm.DB
	
)
