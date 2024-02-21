package rate_limiter

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetRateLimiterConfig(t *testing.T) {
	// Set up test environment
	viper.Set("RATE_LIMIT_DEFAULT_REQUESTS", 5)
	viper.Set("RATE_LIMIT_DEFAULT_EVERY", 30)
	viper.Set("RATE_LIMIT_IP_0", "127.0.0.1")
	viper.Set("RATE_LIMIT_IP_0_REQUESTS", 10)
	viper.Set("RATE_LIMIT_IP_0_EVERY", 60)
	viper.Set("RATE_LIMIT_TOKEN_0", "abc123")
	viper.Set("RATE_LIMIT_TOKEN_0_REQUESTS", 20)
	viper.Set("RATE_LIMIT_TOKEN_0_EVERY", 120)

	expected := RateLimiterConfig{
		Default: Default{
			Requests: 5,
			Every:    30,
		},
		IP: []IP{
			{
				IP:       "127.0.0.1",
				Requests: 10,
				Every:    60,
			},
		},
		Token: []Token{
			{
				Token:    "abc123",
				Requests: 20,
				Every:    120,
			},
		},
	}

	// Call the function under test
	result := GetRateLimiterConfig()

	// Assert the result
	assert.Equal(t, expected, result)
}
