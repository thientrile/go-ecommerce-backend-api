package global

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go-ecommerce-backend-api.com/pkg/logger"
	"go-ecommerce-backend-api.com/pkg/setting"
	"gorm.io/gorm"
)

var (
	Config setting.Config
	Logger *logger.LoggerZap
	RDB    *redis.Client
	MDB    *gorm.DB
	MDBC   *sql.DB
	KafkaProducer *kafka.Writer
)
