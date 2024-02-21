package redis

import (
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string `json:"host,omitempty" env:"REDIS_HOST"`
	Password string `json:"password,omitempty" env:"REDIS_PASSWORD"`
	Port     int    `json:"port,omitempty" env:"REDIS_PORT"`
	DB       int    `json:"db,omitempty" env:"REDIS_DB"`
}

func GetRedisConfig() RedisConfig {
	// set config file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// set default value
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_PORT", 6379)

	// get config
	redisConfig := RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		Port:     viper.GetInt("REDIS_PORT"),
	}

	return redisConfig
}
