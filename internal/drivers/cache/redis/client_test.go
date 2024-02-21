package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	t.Run("Should return the same instance of redis.Client", func(t *testing.T) {
		client1 := GetClient()
		client2 := GetClient()

		assert.Equal(t, client1, client2)
	})
}

func TestConnectRedis(t *testing.T) {
	t.Run("Should return a valid redis.Client instance", func(t *testing.T) {
		client := connectRedis()

		assert.NotNil(t, client)
	})
}
