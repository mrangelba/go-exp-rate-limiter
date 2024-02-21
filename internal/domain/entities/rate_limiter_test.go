package entities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_MarshalBinary(t *testing.T) {
	limiter := RateLimiter{
		Key:       "test_key",
		Requests:  10,
		Every:     60,
		Remaining: 5,
		Reset:     1631234567,
	}

	data, err := limiter.MarshalBinary()
	assert.NoError(t, err)

	var newLimiter RateLimiter
	err = json.Unmarshal(data, &newLimiter)
	assert.NoError(t, err)

	assert.Equal(t, limiter.Key, newLimiter.Key)
	assert.Equal(t, limiter.Requests, newLimiter.Requests)
	assert.Equal(t, limiter.Every, newLimiter.Every)
	assert.Equal(t, limiter.Remaining, newLimiter.Remaining)
	assert.Equal(t, limiter.Reset, newLimiter.Reset)
}

func TestRateLimiter_UnmarshalBinary(t *testing.T) {
	limiter := RateLimiter{
		Key:       "test_key",
		Requests:  10,
		Every:     60,
		Remaining: 5,
		Reset:     1631234567,
	}

	data, err := limiter.MarshalBinary()
	assert.NoError(t, err)

	var newLimiter RateLimiter
	err = newLimiter.UnmarshalBinary(data)
	assert.NoError(t, err)

	assert.Equal(t, limiter.Key, newLimiter.Key)
	assert.Equal(t, limiter.Requests, newLimiter.Requests)
	assert.Equal(t, limiter.Every, newLimiter.Every)
	assert.Equal(t, limiter.Remaining, newLimiter.Remaining)
	assert.Equal(t, limiter.Reset, newLimiter.Reset)
}
