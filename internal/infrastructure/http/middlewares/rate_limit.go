package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases"
)

type rateLimiter struct {
	uc usecases.RateLimitUseCase
}

func NewRateLimiter(uc usecases.RateLimitUseCase) *rateLimiter {
	return &rateLimiter{
		uc: uc,
	}
}

func (m *rateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("API_KEY")

		if apiKey != "" {
			if !m.checkLimitAddHeaders(r.Context(), w, apiKey) {
				return
			}
		} else {
			ips := getIPs(r)

			for _, ip := range ips {
				if !m.checkLimitAddHeaders(r.Context(), w, ip) {
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *rateLimiter) checkLimitAddHeaders(ctx context.Context, w http.ResponseWriter, ip string) bool {
	hasLimit := m.uc.VerifyLimit(ctx, ip)
	headers := m.uc.GetHttpHeaders(ctx, ip)

	for key, value := range headers {
		w.Header().Add(key, value)
	}

	if !hasLimit {
		http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
	}

	return hasLimit
}

func getIPs(r *http.Request) []string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")

		if len(ips) > 0 {
			return ips
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return []string{realIP}
	}

	ips := []string{strings.Split(r.RemoteAddr, ":")[0]}

	return ips
}
