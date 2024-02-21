package config

import (
	"os"
	"testing"

	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/rate_limiter"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/redis"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("Should return the config with default values", func(t *testing.T) {
		config := GetConfig()

		assert.Equal(t, "inmemory", config.Cache)
		assert.NotNil(t, config.Redis)
		assert.NotNil(t, config.RateLimiter)
	})

	t.Run("Should return the config with values from environment variables", func(t *testing.T) {
		// Set environment variables
		_ = viper.BindEnv("CACHE")
		_ = viper.BindEnv("REDIS_HOST")
		_ = viper.BindEnv("REDIS_PORT")
		_ = viper.BindEnv("RATE_LIMIT_DEFAULT_REQUESTS")
		_ = viper.BindEnv("RATE_LIMIT_DEFAULT_EVERY")

		// Set environment variable values
		_ = os.Setenv("CACHE", "redis")
		_ = os.Setenv("REDIS_HOST", "localhost")
		_ = os.Setenv("REDIS_PORT", "6379")
		_ = os.Setenv("RATE_LIMIT_DEFAULT_REQUESTS", "100")
		_ = os.Setenv("RATE_LIMIT_DEFAULT_EVERY", "60")

		config := GetConfig()

		assert.Equal(t, "redis", config.Cache)
		assert.NotNil(t, config.Redis)
		assert.NotNil(t, config.RateLimiter)
		assert.Equal(t, "localhost", config.Redis.Host)
		assert.Equal(t, 6379, config.Redis.Port)
		assert.Equal(t, 100, config.RateLimiter.Default.Requests)
		assert.Equal(t, 60, config.RateLimiter.Default.Every)
	})
}

func TestConfig_String(t *testing.T) {
	t.Run("Should return the JSON representation of the config", func(t *testing.T) {
		config := Config{
			Cache: "inmemory",
			Redis: redis.RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
			RateLimiter: rate_limiter.RateLimiterConfig{
				Default: rate_limiter.Default{
					Requests: 10,
					Every:    60,
				},
			},
		}

		expectedJSON := `{"cache":"inmemory","redis":{"host":"localhost","port":6379},"rate_limiter":{"default":{"requests":10,"every":60}}}`

		assert.Equal(t, expectedJSON, config.String())
	})
}
