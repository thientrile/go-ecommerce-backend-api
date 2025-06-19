package initialize

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go-ecommerce-backend-api.com/global"
	"go.uber.org/zap"
)

var ctx = context.Background()


func InitRedis() {
	// Initialize Redis connection here	
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:	fmt.Sprintf("%s:%v",r.Host,r.Port), // Redis server address
		Password: r.Password,            // No password set
		DB:       r.Database,             // Use default DB
		PoolSize: r.Pool_size, // Set pool size
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Panic("Redis Initialzation Error::", zap.Error(err))
	}
	global.Logger.Info("✅ Redis connection pool initialized successfully")
	fmt.Println("Initialized Redis successfully");
	global.RDB = rdb

}