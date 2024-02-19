package usecases

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain"
	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/entities"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config/rate_limiter"
)

type RateLimitUseCase interface {
	GetHttpHeaders(ctx context.Context, key string) map[string]string
	VerifyLimit(ctx context.Context, key string) bool
}

type rateLimitUseCase struct {
	config config.Config
	cache  domain.RateLimitCache
}

func NewRateLimitUseCase(config config.Config, cache domain.RateLimitCache) RateLimitUseCase {
	return &rateLimitUseCase{
		config: config,
		cache:  cache,
	}
}

func (uc *rateLimitUseCase) VerifyLimit(ctx context.Context, key string) bool {
	_, err := uc.cache.Get(ctx, key)

	if err != nil {
		every := uc.config.RateLimiter.Default.Every
		requests := uc.config.RateLimiter.Default.Requests

		if slices.ContainsFunc(uc.config.RateLimiter.Token, func(s rate_limiter.Token) bool {
			return s.Token == key
		}) {
			index := slices.IndexFunc(uc.config.RateLimiter.Token, func(s rate_limiter.Token) bool {
				return s.Token == key
			})

			every = uc.config.RateLimiter.Token[index].Every
			requests = uc.config.RateLimiter.Token[index].Requests
		} else if slices.ContainsFunc(uc.config.RateLimiter.IP, func(s rate_limiter.IP) bool {
			return s.IP == key
		}) {
			index := slices.IndexFunc(uc.config.RateLimiter.IP, func(s rate_limiter.IP) bool {
				return s.IP == key
			})

			every = uc.config.RateLimiter.IP[index].Every
			requests = uc.config.RateLimiter.IP[index].Requests
		}

		uc.cache.Set(ctx,
			entities.RateLimiter{
				Key:       key,
				Every:     every,
				Remaining: requests,
				Requests:  0,
				Reset:     time.Now().Add(time.Duration(every) * time.Second).Unix(),
			}, time.Duration(every)*time.Second)

		return uc.VerifyLimit(ctx, key)
	}

	limit, err := uc.validateCacheLimit(ctx, key)

	if err != nil {
		log.Println(err)
		return false
	}

	return limit
}

func (uc *rateLimitUseCase) validateCacheLimit(ctx context.Context, key string) (bool, error) {
	rate, err := uc.cache.Get(ctx, key)
	if err != nil {
		return false, err
	}

	if rate.Remaining <= 0 && rate.Every > 0 {
		return false, nil
	}

	rate.Requests++

	if rate.Every > 0 && rate.Remaining > 0 {
		rate.Remaining--
	}

	every := (time.Duration(rate.Reset) - time.Duration(time.Now().Unix())) * time.Second

	if err := uc.cache.Set(ctx, *rate, every); err != nil {
		return false, err
	}

	return true, nil
}

func (uc *rateLimitUseCase) GetHttpHeaders(ctx context.Context, key string) map[string]string {
	rate, err := uc.cache.Get(ctx, key)
	if err != nil {
		return map[string]string{}
	}

	every := (time.Duration(rate.Reset) - time.Duration(time.Now().Unix())) * time.Second

	headers := map[string]string{
		"Ratelimit-Limit":     fmt.Sprintf("%v", rate.Requests),
		"Ratelimit-Remaining": fmt.Sprintf("%v", rate.Remaining),
		"Ratelimit-Reset":     fmt.Sprintf("%v", every),
	}

	return headers
}
