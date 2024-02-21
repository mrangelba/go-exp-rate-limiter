package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCacheStrategy(t *testing.T) {
	t.Run("Should return Redis cache strategy when cache is set to 'redis'", func(t *testing.T) {
		cache := "redis"
		expected := NewRateLimitRedis()

		result := GetCacheStrategy(cache)

		assert.Equal(t, expected, result)
	})

	t.Run("Should return InMemory cache strategy when cache is set to any value other than 'redis'", func(t *testing.T) {
		cache := "memcached"
		expected := NewRateLimitInMemory()

		result := GetCacheStrategy(cache)

		assert.Equal(t, expected, result)
	})
}
