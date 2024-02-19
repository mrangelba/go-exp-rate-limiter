package strategies

import (
	"context"
	"errors"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain"
	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
)

type rateLimitInMemory struct {
	rates map[string]entities.RateLimiter
}

func NewRateLimitInMemory() domain.RateLimitCache {
	return &rateLimitInMemory{
		rates: make(map[string]entities.RateLimiter),
	}
}

func (r *rateLimitInMemory) Set(ctx context.Context, rate entities.RateLimiter, every time.Duration) error {
	r.rates[rate.Key] = rate
	return nil
}

func (r *rateLimitInMemory) Get(ctx context.Context, key string) (*entities.RateLimiter, error) {
	rate, ok := r.rates[key]
	if !ok {
		return nil, errors.New("rate limit not found")
	}

	if rate.Reset < time.Now().Unix() {
		delete(r.rates, key)
		return nil, errors.New("rate limit expired")
	}

	return &rate, nil
}
