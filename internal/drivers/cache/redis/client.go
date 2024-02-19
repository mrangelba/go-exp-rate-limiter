package redis

import (
	"fmt"

	"sync"

	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/redis/go-redis/v9"
)

var once sync.Once
var instance *redis.Client

func GetClient() *redis.Client {
	once.Do(func() {
		instance = connectRedis()
	})

	return instance
}

func connectRedis() *redis.Client {
	cfg := config.GetConfig().Redis

	client := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})

	return client
}
