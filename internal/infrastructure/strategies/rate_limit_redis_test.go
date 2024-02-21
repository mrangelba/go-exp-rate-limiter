package strategies

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRateLimitRedis_Set(t *testing.T) {
	ctx := context.TODO()
	rate := entities.RateLimiter{
		Key:      "test_key",
		Requests: 10,
	}

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}

	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	t.Run("Should set rate limit in Redis", func(t *testing.T) {
		endpoint, err := redisC.Endpoint(ctx, "")
		assert.NoError(t, err)

		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: endpoint,
		})

		// Create a rateLimitRedis instance
		rl := rateLimitRedis{
			client: client,
		}

		// Set the rate limit
		err = rl.Set(ctx, rate, time.Minute)

		// Check if there was no error
		assert.NoError(t, err)

		// Get the rate limit from Redis
		val, err := client.Get(ctx, rate.Key).Result()

		// Check if there was no error
		assert.NoError(t, err)

		// Unmarshal the rate limit from JSON
		var storedRate entities.RateLimiter
		err = json.Unmarshal([]byte(val), &storedRate)

		// Check if there was no error
		assert.NoError(t, err)

		// Check if the stored rate limit matches the original rate limit
		assert.Equal(t, rate, storedRate)
	})

	t.Run("Should handle error when setting rate limit in Redis", func(t *testing.T) {
		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: "",
		})

		client.FlushAll(ctx)

		// Create a rateLimitRedis instance
		rl := rateLimitRedis{
			client: client,
		}

		// Set the rate limit
		err = rl.Set(ctx, rate, time.Minute)

		// Check if there was no error
		assert.Error(t, err)
	})
}

func TestRateLimitRedis_Get(t *testing.T) {
	ctx := context.TODO()
	key := "test_key"
	expectedRate := entities.RateLimiter{
		Key:      key,
		Requests: 10,
	}

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

	t.Run("Should get rate limit from Redis", func(t *testing.T) {
		endpoint, err := redisC.Endpoint(ctx, "")
		assert.NoError(t, err)

		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: endpoint,
		})

		// Create a rateLimitRedis instance
		rl := rateLimitRedis{
			client: client,
		}

		// Set the rate limit in Redis
		err = client.Set(ctx, key, `{"key":"test_key","requests":10}`, 0).Err()
		assert.NoError(t, err)

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		// Check if there was no error
		assert.NoError(t, err)

		// Check if the retrieved rate limit matches the expected rate limit
		assert.Equal(t, &expectedRate, rate)
	})

	t.Run("Should handle error when getting rate limit from Redis", func(t *testing.T) {
		endpoint, err := redisC.Endpoint(ctx, "")
		assert.NoError(t, err)

		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: endpoint,
		})

		// Create a rateLimitRedis instance
		rl := rateLimitRedis{
			client: client,
		}

		// Set the rate limit in Redis
		err = client.Set(ctx, key, `{"key":"test_key","requests":10.0}`, 0).Err()
		assert.NoError(t, err)

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		// Check if the error is not nil
		assert.Error(t, err)

		// Check if the retrieved rate limit is nil
		assert.Nil(t, rate)
	})

	t.Run("Should handle error when rate limit does not exist in Redis", func(t *testing.T) {
		endpoint, err := redisC.Endpoint(ctx, "")
		assert.NoError(t, err)

		// Create a Redis client
		client := redis.NewClient(&redis.Options{
			Addr: endpoint,
		})
		client.FlushAll(ctx)

		// Create a rateLimitRedis instance
		rl := rateLimitRedis{
			client: client,
		}

		// Get the rate limit
		rate, err := rl.Get(ctx, key)

		// Check if the error is not nil
		assert.Error(t, err)

		// Check if the retrieved rate limit is nil
		assert.Nil(t, rate)
	})
}
