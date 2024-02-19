package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrangelba/go-exp-rate-limiter/internal/domain/usecases"
	"github.com/mrangelba/go-exp-rate-limiter/internal/drivers/config"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/http/middlewares"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/strategies"
	"github.com/mrangelba/go-exp-rate-limiter/internal/infrastructure/utils"
)

func main() {
	log.Println("Starting server...")

	config := config.GetConfig()

	uc := usecases.NewRateLimitUseCase(
		config,
		strategies.GetCacheStrategy(config.Cache),
	)

	rateLimit := middlewares.NewRateLimiter(uc)

	router := gin.Default()
	router.Use(
		utils.MiddlewareToGin(rateLimit.Handler),
	)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Go Expert!")
	})

	router.Run(":8080")
}
