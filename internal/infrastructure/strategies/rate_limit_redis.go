package strategies

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain"
	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"

	"github.com/redis/go-redis/v9"

	driversRedis "github.com/mrangelba/go-exp-rate-limiter/internal/drivers/cache/redis"
)

type rateLimitRedis struct {
	client *redis.Client
}

func NewRateLimitRedis() domain.RateLimitCache {
	return &rateLimitRedis{
		client: driversRedis.GetClient(),
	}
}

func (r *rateLimitRedis) Set(ctx context.Context, rate entities.RateLimiter, every time.Duration) error {
	json, err := json.Marshal(rate)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, rate.Key, json, every).Err()
}

func (r *rateLimitRedis) Get(ctx context.Context, key string) (*entities.RateLimiter, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Error reading rate limit for key:", key)
		}

		return nil, err
	}

	var rate entities.RateLimiter
	if err := json.Unmarshal([]byte(val), &rate); err != nil {
		log.Println("Error unmarshalling rate limit for key:", key)

		return nil, err
	}

	return &rate, nil
}
