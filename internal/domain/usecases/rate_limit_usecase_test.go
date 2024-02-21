package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/rate_limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRateLimitCache struct {
	mock.Mock
}

func (m *mockRateLimitCache) Get(ctx context.Context, key string) (*entities.RateLimiter, error) {
	args := m.Called(ctx, key)

	if args.Get(0) != nil {
		return args.Get(0).(*entities.RateLimiter), nil
	}

	return nil, args.Error(1)
}

func (m *mockRateLimitCache) Set(ctx context.Context, rate entities.RateLimiter, duration time.Duration) error {
	args := m.Called(ctx, rate, duration)
	return args.Error(0)
}

type mockRateLimitCache2 struct {
	mock.Mock
	getCallCount int
}

func (m *mockRateLimitCache2) Get(ctx context.Context, key string) (*entities.RateLimiter, error) {
	args := m.Called(ctx, key)

	if args.Get(0) != nil && m.getCallCount == 1 {
		return args.Get(0).(*entities.RateLimiter), nil
	}

	m.getCallCount++

	return nil, args.Error(1)
}

func (m *mockRateLimitCache2) Set(ctx context.Context, rate entities.RateLimiter, duration time.Duration) error {
	args := m.Called(ctx, rate, duration)
	return args.Error(0)
}

func TestRateLimitUseCase_VerifyLimit(t *testing.T) {
	config := config.Config{
		RateLimiter: rate_limiter.RateLimiterConfig{
			Default: rate_limiter.Default{
				Every:    60,
				Requests: 100,
			},
			Token: []rate_limiter.Token{
				{
					Token:    "token1",
					Every:    30,
					Requests: 50,
				},
			},
			IP: []rate_limiter.IP{
				{
					IP:       "127.0.0.1",
					Every:    10,
					Requests: 20,
				},
			},
		},
	}

	t.Run("Should return true when rate limit is not exceeded", func(t *testing.T) {
		ctx := context.Background()
		key := "token1"
		cache := new(mockRateLimitCache)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       "token1",
			Every:     30,
			Remaining: 10,
			Requests:  40,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, nil)
		cache.On("Set", ctx, mock.Anything, mock.Anything).Return(nil).Once()

		result := useCase.VerifyLimit(ctx, key)

		assert.True(t, result)
	})

	t.Run("Should return false when set cache error", func(t *testing.T) {
		ctx := context.Background()
		key := "token1"
		cache := new(mockRateLimitCache)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       "token1",
			Every:     30,
			Remaining: 10,
			Requests:  40,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, nil)
		cache.On("Set", ctx, mock.Anything, mock.Anything).Return(errors.New("error")).Once()

		result := useCase.VerifyLimit(ctx, key)

		assert.False(t, result)
	})

	t.Run("Should return false when rate limit is exceeded", func(t *testing.T) {
		ctx := context.Background()
		key := "token1"

		cache := new(mockRateLimitCache)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       key,
			Every:     30,
			Remaining: 0,
			Requests:  50,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, nil)

		result := useCase.VerifyLimit(ctx, key)

		assert.False(t, result)
	})
}

func TestRateLimitUseCase_VerifyLimit_ReadConfig(t *testing.T) {
	config := config.Config{
		RateLimiter: rate_limiter.RateLimiterConfig{
			Default: rate_limiter.Default{
				Every:    60,
				Requests: 100,
			},
			Token: []rate_limiter.Token{
				{
					Token:    "token1",
					Every:    30,
					Requests: 50,
				},
			},
			IP: []rate_limiter.IP{
				{
					IP:       "127.0.0.1",
					Every:    10,
					Requests: 20,
				},
			},
		},
	}

	t.Run("Should return true when rate limit is not exceeded Read config token", func(t *testing.T) {
		ctx := context.Background()
		key := "token1"
		cache := new(mockRateLimitCache2)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       "token1",
			Every:     30,
			Remaining: 10,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, errors.New("not found")).Times(3)
		cache.On("Set", ctx, mock.Anything, mock.Anything).Return(nil)

		result := useCase.VerifyLimit(ctx, key)

		assert.True(t, result)
	})

	t.Run("Should return true when rate limit is not exceeded Read config IP", func(t *testing.T) {
		ctx := context.Background()
		key := "127.0.0.1"
		cache := new(mockRateLimitCache2)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       "127.0.0.1",
			Every:     30,
			Remaining: 10,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, errors.New("not found")).Times(3)
		cache.On("Set", ctx, mock.Anything, mock.Anything).Return(nil)

		result := useCase.VerifyLimit(ctx, key)

		assert.True(t, result)
	})
}

func TestRateLimitUseCase_GetHttpHeaders(t *testing.T) {
	config := config.Config{
		RateLimiter: rate_limiter.RateLimiterConfig{
			Default: rate_limiter.Default{
				Every:    60,
				Requests: 100,
			},
			Token: []rate_limiter.Token{
				{
					Token:    "token1",
					Every:    30,
					Requests: 50,
				},
			},
			IP: []rate_limiter.IP{
				{
					IP:       "127.0.0.1",
					Every:    10,
					Requests: 20,
				},
			},
		},
	}

	t.Run("Should return HTTP headers with rate limit information", func(t *testing.T) {
		ctx := context.Background()
		key := "token1"

		cache := new(mockRateLimitCache)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(&entities.RateLimiter{
			Key:       key,
			Every:     30,
			Remaining: 10,
			Requests:  40,
			Reset:     time.Now().Add(30 * time.Second).Unix(),
		}, nil)

		headers := useCase.GetHttpHeaders(ctx, key)

		assert.Equal(t, "40", headers["Ratelimit-Limit"])
		assert.Equal(t, "10", headers["Ratelimit-Remaining"])
		assert.Equal(t, "30s", headers["Ratelimit-Reset"])
	})

	t.Run("Should return empty HTTP headers when rate limit is not found", func(t *testing.T) {
		ctx := context.Background()
		key := "token2"

		cache := new(mockRateLimitCache)

		useCase := usecases.NewRateLimitUseCase(config, cache)

		cache.On("Get", ctx, key).Return(nil, errors.New("not found"))

		headers := useCase.GetHttpHeaders(ctx, key)

		assert.Empty(t, headers)
	})
}
