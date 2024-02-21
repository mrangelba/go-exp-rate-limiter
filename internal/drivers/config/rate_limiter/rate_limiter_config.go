package rate_limiter

import (
	"fmt"

	"github.com/spf13/viper"
)

type RateLimiterConfig struct {
	Default Default `json:"default"`
	IP      []IP    `json:"ip,omitempty"`
	Token   []Token `json:"token,omitempty"`
}

type Default struct {
	Requests int `json:"requests,omitempty"`
	Every    int `json:"every,omitempty"`
}

type IP struct {
	IP       string `json:"ip,omitempty"`
	Requests int    `json:"requests,omitempty"`
	Every    int    `json:"every,omitempty"`
}

type Token struct {
	Token    string `json:"token,omitempty"`
	Requests int    `json:"requests,omitempty"`
	Every    int    `json:"every,omitempty"`
}

// GetRateLimiterConfig returns the rate limiter configuration
func GetRateLimiterConfig() RateLimiterConfig {
	// set config file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// set default value
	viper.SetDefault("RATE_LIMIT_DEFAULT_REQUESTS", 10)
	viper.SetDefault("RATE_LIMIT_DEFAULT_EVERY", 60)

	// get config default
	rateLimiterConfig := RateLimiterConfig{
		Default: Default{
			Requests: viper.GetInt("RATE_LIMIT_DEFAULT_REQUESTS"),
			Every:    viper.GetInt("RATE_LIMIT_DEFAULT_EVERY"),
		},
	}

	for i := 0; ; i++ {
		ipKey := fmt.Sprintf("RATE_LIMIT_IP_%d", i)

		if !viper.IsSet(ipKey) {
			break
		}

		ip := viper.GetString(ipKey)
		requests := viper.GetInt(fmt.Sprintf("RATE_LIMIT_IP_%d_REQUESTS", i))
		every := viper.GetInt(fmt.Sprintf("RATE_LIMIT_IP_%d_EVERY", i))

		rateLimiterConfig.IP = append(rateLimiterConfig.IP, IP{
			IP:       ip,
			Requests: requests,
			Every:    every,
		})
	}

	for i := 0; ; i++ {
		tokenKey := fmt.Sprintf("RATE_LIMIT_TOKEN_%d", i)

		if !viper.IsSet(tokenKey) {
			break
		}

		token := viper.GetString(tokenKey)
		requests := viper.GetInt(fmt.Sprintf("RATE_LIMIT_TOKEN_%d_REQUESTS", i))
		every := viper.GetInt(fmt.Sprintf("RATE_LIMIT_TOKEN_%d_EVERY", i))

		rateLimiterConfig.Token = append(rateLimiterConfig.Token, Token{
			Token:    token,
			Requests: requests,
			Every:    every,
		})
	}

	return rateLimiterConfig
}
