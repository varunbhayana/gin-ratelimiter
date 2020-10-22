package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/varunbhayana/gin-ratelimiter/service"
	"github.com/varunbhayana/gin-ratelimiter/util/cycle_util"
)

type RedisModel struct {
	Time  int64 `json:"Time"`
	Count int64 `json:"Count"`
}

func RateLimit() func(*gin.Context) {
	return func(c *gin.Context) {
		cycle_util.DegdCall(
			1*time.Second,
			c,
			func() (int, interface{}) {
				userId := c.GetHeader("user-id")
				if userId == "" {
					return 400, "Bad Request"
				}

				return service.RateLimiter.RateLimit(userId)
			},
		)
	}

}
