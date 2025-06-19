package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go-ecommerce-backend-api.com/global"
)

func GetCache(cxt context.Context, key string, obj interface{}) error {
	// This function should implement the logic to get a value from the cache
	// using the provided context and key.
	// For now, we return an empty string and nil error as a placeholder.
	rs, err := global.RDB.Get(cxt, key).Result()

	if err == redis.Nil {
		return fmt.Errorf("key %s does not exist", key)
	}
	if err != nil {
		return fmt.Errorf("error getting key %s: %v", key, err)
	}
	// convert json string to obj
	if err := json.Unmarshal([]byte(rs), obj); err != nil {
		return fmt.Errorf("error unmarshalling json for key %s: %v", key, err)
	}
	return nil
}
