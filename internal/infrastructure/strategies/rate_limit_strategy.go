package strategies

import (
	"log"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain"
)

func GetCacheStrategy(cache string) domain.RateLimitCache {
	if cache == "redis" {
		log.Println("Using Redis as cache")
		return NewRateLimitRedis()
	}

	log.Println("Using InMemory as cache")
	return NewRateLimitInMemory()
}
