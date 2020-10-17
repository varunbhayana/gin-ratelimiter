package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_redis "github.com/go-redis/redis"

	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/varunbhayana/gin-ratelimiter/db"
	"github.com/varunbhayana/gin-ratelimiter/service"
)

type RedisModel struct {
	Time  int64 `json:"Time"`
	Count int64 `json:"Count"`
}

func RateLimit(c *gin.Context) {

	now := time.Now()
	userId := c.GetHeader("user-id")
	client := db.GetRedis()

	locker := service.NewRedisLock(client)

	lock, err := locker.Obtain(userId+"_lock", 100*time.Millisecond, &service.Options{
		Context:       client.Context(),
		Metadata:      "_lock",
		RetryStrategy: service.LinearBackoff(60 * time.Millisecond),
	})
	if err == redislock.ErrNotObtained {
		c.String(408, "Retry")
		return
	} else if err != nil {
		c.String(408, "Retry")
		return
	}
	defer lock.Release()
	fmt.Println("I have a lock!")
	value, err := client.Get(userId).Result()
	if err == nil {
		var redisValue map[int64]int64
		if err := json.Unmarshal([]byte(value), &redisValue); err != nil {
			panic(err)
		}
		if validate(&redisValue, now) {
			c.String(http.StatusOK, "ok")
			setForUser(&redisValue, userId, now, client)
		} else {
			// too many request response
			c.String(429, "too many requests")
		}
		return
	} else {
		setForUser(&map[int64]int64{}, userId, now, client)
		c.String(200, "ok")
		return
	}

}

func validate(user *map[int64]int64, now time.Time) bool {

	after := now.Add(-1*time.Hour).Unix() / 60

	total := int64(0)
	for k, v := range *user {
		if k < after {
			delete(*user, k)
		} else {
			total += v
		}
	}
	// in 1 hour max 1k request
	if total+1 > 1000 {
		return false
	}

	minute := now.Unix() / 60
	if val, ok := (*user)[minute]; ok {
		if val >= 3 {
			return false
		} else {
			return true
		}
	}
	return true
}

func setForUser(redisValue *map[int64]int64, userId string, now time.Time, client *_redis.Client) {
	key := now.Unix() / 60
	if val, ok := (*redisValue)[key]; ok {
		(*redisValue)[key] = val + 1

	} else {
		(*redisValue)[key] = 1
	}

	b, _ := json.Marshal(redisValue)

	ok, err := client.Set(userId, string(b), time.Duration(1*time.Hour)).Result()
	fmt.Println(ok, err)
}
