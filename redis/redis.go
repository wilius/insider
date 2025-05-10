package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"insider/configs"
	"sync"
)

var (
	instance       *redis.Client
	connectionOnce sync.Once
)

func Instance() *redis.Client {
	connectionOnce.Do(func() {
		redisConfig := configs.Instance().
			GetRedis()
		instance = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", redisConfig.GetHost(), redisConfig.GetPort()),
		})
	})

	return instance
}
