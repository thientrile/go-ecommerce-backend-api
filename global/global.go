package global

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/ulule/limiter/v3"
	"go-ecommerce-backend-api.com/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go-ecommerce-backend-api.com/pkg/setting"
	"gorm.io/gorm"
)

var (
	Config         setting.Config
	Logger         *logger.LoggerZap
	RDB            *redis.Client
	MDB            *gorm.DB
	MDBC           *sql.DB
	KafkaProducers map[string]*kafka.Writer
	Limiters       map[string]*limiter.Limiter
	OpsProcessed   prometheus.Counter
)
