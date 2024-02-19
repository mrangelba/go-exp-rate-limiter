package domain

import (
	"context"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
)

type RateLimitCache interface {
	Set(ctx context.Context, rate entities.RateLimiter, every time.Duration) error
	Get(ctx context.Context, key string) (*entities.RateLimiter, error)
}
