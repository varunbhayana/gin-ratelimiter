package handler

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	_redis "github.com/go-redis/redis"
	"github.com/varunbhayana/gin-ratelimiter/db"
	"github.com/varunbhayana/gin-ratelimiter/service"
	"github.com/varunbhayana/gin-ratelimiter/util/cycle_util"
)

var MAX_MINUTE int

var MAX_HOUR int

func init() {
	MAX_MINUTE, _ = strconv.Atoi(os.Getenv("MAX_MINUTE"))
	MAX_HOUR, _ = strconv.Atoi(os.Getenv("MAX_HOUR"))

}

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
					return 408, "Retry"
				} else if err != nil {
					return 408, "Retry"

				}
				defer func() {
					err := lock.Release()
					if err != nil {
						if x := recover(); x == nil {
							panic("retry")
						}

					}
				}()
				value, err := client.Get(userId).Result()
				if err == nil {
					var redisValue map[int64]int
					if err := json.Unmarshal([]byte(value), &redisValue); err != nil {
						panic(err)
					}
					if validate(&redisValue, now) {
						setForUser(&redisValue, userId, now, client)
						return 200, "ok"
					} else {
						// too many request response
						return 429, "too many requests"
					}
				} else {
					setForUser(&map[int64]int{}, userId, now, client)
					return 200, "ok"

				}
			},
		)
	}

}

func validate(user *map[int64]int, now time.Time) bool {

	after := now.Add(-1*time.Hour).Unix() / 60

	total := 0
	for k, v := range *user {
		if k < after {
			delete(*user, k)
		} else {
			total += v
		}
	}
	// in 1 hour max 1k request
	if total+1 > MAX_HOUR {
		return false
	}

	minute := now.Unix() / 60
	if val, ok := (*user)[minute]; ok {
		if val >= MAX_MINUTE {
			return false
		} else {
			return true
		}
	}
	return true
}

func setForUser(redisValue *map[int64]int, userId string, now time.Time, client *_redis.Client) {
	key := now.Unix() / 60
	if val, ok := (*redisValue)[key]; ok {
		(*redisValue)[key] = val + 1

	} else {
		(*redisValue)[key] = 1
	}

	b, _ := json.Marshal(redisValue)

	client.Set(userId, string(b), time.Duration(1*time.Hour)).Result()
}
