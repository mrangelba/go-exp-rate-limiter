package strategies

import (
	"context"
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
	err := r.client.Set(ctx, rate.Key, rate, every).Err()

	if err != nil {
		log.Println(err)
	}

	return err

}

func (r *rateLimitRedis) Get(ctx context.Context, key string) (*entities.RateLimiter, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var rate entities.RateLimiter
	err = rate.UnmarshalBinary([]byte(val))
	if err != nil {
		return nil, err
	}

	return &rate, nil
}
