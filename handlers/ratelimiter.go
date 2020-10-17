package handler

import (
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/varunbhayana/gin-ratelimiter/db"
	"github.com/varunbhayana/gin-ratelimiter/service"
)

type RedisModel struct {
	time  int64
	count int64
}

func RateLimit(c *gin.Context) {

	userId := c.GetHeader("user-id")
	client := db.GetRedis()

	locker := service.NewRedisLock(client)

	// Try to obtain lock.
	lock, err := locker.Obtain(userId, 200*time.Millisecond, nil)
	if err == redislock.ErrNotObtained {
		c.String(408, "Retry")
	} else if err != nil {
		c.String(408, "Retry")
	}
	defer lock.Release()
	fmt.Println("I have a lock!")
	value, err := client.Get(userId).Result()
	if err == nil {
		fmt.Println(value)

	} else {
		redisValue := make([]*RedisModel, 0)

		redisValue = append(redisValue, &RedisModel{
			time:  time.Now().Unix() / 60,
			count: 1,
		})

		client.Set(userId, redisValue, time.Duration(1*time.Hour)).Result()
	}

}
