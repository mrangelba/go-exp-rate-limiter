package redis

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRedisConfig(t *testing.T) {
	// Set up test environment
	viper.Reset()
	viper.Set("REDIS_HOST", "testhost")
	viper.Set("REDIS_PASSWORD", "testpassword")
	viper.Set("REDIS_DB", 1)
	viper.Set("REDIS_PORT", 1234)

	expected := RedisConfig{
		Host:     "testhost",
		Password: "testpassword",
		DB:       1,
		Port:     1234,
	}

	// Call the function under test
	result := GetRedisConfig()

	// Assert the result
	assert.Equal(t, expected, result)
}
