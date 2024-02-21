package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRateLimitUseCase struct{}

func (m *mockRateLimitUseCase) VerifyLimit(ctx context.Context, key string) bool {
	// Mock implementation
	return true
}

func (m *mockRateLimitUseCase) GetHttpHeaders(ctx context.Context, key string) map[string]string {
	// Mock implementation
	return map[string]string{
		"X-RateLimit-Limit":     "100",
		"X-RateLimit-Remaining": "50",
	}
}

type mockRateLimitUseCaseError struct{}

func (m *mockRateLimitUseCaseError) VerifyLimit(ctx context.Context, key string) bool {
	// Mock implementation
	return false
}

func (m *mockRateLimitUseCaseError) GetHttpHeaders(ctx context.Context, key string) map[string]string {
	// Mock implementation
	return map[string]string{
		"X-RateLimit-Limit":     "100",
		"X-RateLimit-Remaining": "0",
	}
}

func TestRateLimiter_Handler(t *testing.T) {
	uc := &mockRateLimitUseCase{}
	rl := NewRateLimiter(uc)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Should call next handler when API_KEY is provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)
		req.Header.Set("API_KEY", "test_key")

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should return 429 Too Many Requests when rate limit is reached", func(t *testing.T) {
		uc := new(mockRateLimitUseCaseError)
		rl := NewRateLimiter(uc)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("API_KEY", "test_key")
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", rr.Body.String())
	})

	t.Run("Should call next handler when X-Forwarded-For header is provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)
		req.Header.Set("X-Forwarded-For", "127.0.0.1, 192.168.0.1")

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should call next handler when X-Real-IP header is provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)
		req.Header.Set("X-Real-IP", "127.0.0.1")

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should call next handler when no API_KEY or IP headers are provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Should return 429 Too Many Requests when rate limit is reached", func(t *testing.T) {
		uc := new(mockRateLimitUseCaseError)
		rl := NewRateLimiter(uc)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", rr.Body.String())
	})

	t.Run("Should add rate limit headers to the response", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		rl.Handler(handler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "100", rr.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "50", rr.Header().Get("X-RateLimit-Remaining"))
	})
}
