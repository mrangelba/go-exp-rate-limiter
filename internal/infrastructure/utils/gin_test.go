package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareToGin(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create a test route that uses the MiddlewareToGin function
	router.GET("/test", MiddlewareToGin(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Perform some middleware logic here
			// For example, set a custom header
			w.Header().Set("X-Custom-Header", "Test")
			next.ServeHTTP(w, r)
		})
	}))

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Assert that the middleware logic was executed
	assert.Equal(t, "Test", recorder.Header().Get("X-Custom-Header"))
}

func TestMiddlewareToGinError(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create a test route that uses the MiddlewareToGin function
	router.GET("/test", MiddlewareToGin(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Perform some middleware logic here
			// For example, set a custom header
			http.Error(w, "Test", http.StatusInternalServerError)
		})
	}))

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create a test response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Assert that the middleware logic was executed
	assert.Equal(t, http.StatusInternalServerError, recorder.Result().StatusCode)
}
