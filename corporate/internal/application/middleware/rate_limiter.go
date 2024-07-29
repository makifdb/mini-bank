package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

func RateLimiterMiddleware(rps int) gin.HandlerFunc {
	rl := ratelimit.New(rps)

	return func(c *gin.Context) {
		rl.Take()
		c.Next()
	}
}
