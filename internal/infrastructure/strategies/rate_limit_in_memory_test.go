package strategies

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitInMemory_Set(t *testing.T) {
	ctx := context.TODO()
	rate := entities.RateLimiter{
		Key:      "test_key",
		Requests: 10,
	}

	t.Run("Should set rate limit in memory", func(t *testing.T) {
		// Create a rateLimitInMemory instance
		rl := NewRateLimitInMemory()

		// Set the rate limit
		err := rl.Set(ctx, rate, time.Minute)

		// Check if there was no error
		assert.NoError(t, err)
	})
}

func TestRateLimitInMemory_Get(t *testing.T) {
	ctx := context.TODO()
	key := "test_key"
	expectedRate := entities.RateLimiter{
		Key:      key,
		Requests: 10,
		Reset:    int64(time.Minute),
	}

	t.Run("Should get rate limit from memory", func(t *testing.T) {
		// Create a rateLimitInMemory instance
		rl := NewRateLimitInMemory()

		// Set the rate limit in memory
		err := rl.Set(ctx, expectedRate, time.Minute)
		assert.NoError(t, err)

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		log.Println(rate)

		// Check if there was no error
		assert.NoError(t, err)

		// Check if the retrieved rate limit matches the expected rate limit
		assert.Equal(t, &expectedRate, rate)
	})

	t.Run("Should handle error when rate limit does not exist in memory", func(t *testing.T) {
		// Create a rateLimitInMemory instance
		rl := NewRateLimitInMemory()

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		// Check if the error is not nil
		assert.Error(t, err)
		assert.ErrorContains(t, err, "rate limit not found")

		// Check if the retrieved rate limit is nil
		assert.Nil(t, rate)
	})

	t.Run("Should handle error when rate limit has expired", func(t *testing.T) {
		// Create a rateLimitInMemory instance
		rl := NewRateLimitInMemory()

		// Set the rate limit in memory with an expired reset time
		expiredRate := entities.RateLimiter{
			Key:      key,
			Requests: 10,
			Reset:    time.Now().Unix() - 1,
		}
		err := rl.Set(ctx, expiredRate, time.Minute)
		assert.NoError(t, err)

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		// Check if the error is not nil
		assert.Error(t, err)
		assert.ErrorContains(t, err, "rate limit expired")

		// Check if the retrieved rate limit is nil
		assert.Nil(t, rate)
	})
}
