package config

import (
	"encoding/json"

	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/rate_limiter"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/redis"
	"github.com/spf13/viper"
)

type Config struct {
	Cache       string                         `json:"cache"`
	Redis       redis.RedisConfig              `json:"redis"`
	RateLimiter rate_limiter.RateLimiterConfig `json:"rate_limiter"`
}

func GetConfig() Config {
	// set config file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// set default value
	viper.SetDefault("CACHE", "inmemory")

	return Config{
		Cache:       viper.GetString("CACHE"),
		Redis:       redis.GetRedisConfig(),
		RateLimiter: rate_limiter.GetRateLimiterConfig(),
	}
}

func (r Config) String() string {
	data, _ := json.Marshal(r)

	return string(data)
}
