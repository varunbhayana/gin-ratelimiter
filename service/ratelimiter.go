package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bsm/redislock"
	_redis "github.com/go-redis/redis"
	"github.com/varunbhayana/rate-limiting/db"
	"github.com/varunbhayana/rate-limiting/enum"
)

type ratelimiter struct{}

var (
	RateLimiter = ratelimiter{}
)
var MAX_MINUTE int

var MAX_HOUR int

func init() {
	MAX_MINUTE, _ = strconv.Atoi(enum.ReadEnv("MAX_MINUTE"))
	MAX_HOUR, _ = strconv.Atoi(enum.ReadEnv("MAX_HOUR"))

}

func (limiter *ratelimiter) RateLimit(userId string, applicationId string) (int, string) {
	now := time.Now()
	client := db.GetRedis()

	locker := NewRedisLock(client)

	lock, err := locker.Obtain(getLockKey(userId, applicationId), 100*time.Millisecond, &Options{
		Context:       client.Context(),
		Metadata:      "_lock",
		RetryStrategy: LinearBackoff(60 * time.Millisecond),
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
	value, err := client.Get(getRateLimitingKey(userId, applicationId)).Result()
	if err == nil {
		var redisValue map[int64]int
		if err := json.Unmarshal([]byte(value), &redisValue); err != nil {
			panic(err)
		}
		if rateInLimit(&redisValue, now) {
			incrementRequestCount(&redisValue, userId, applicationId, now, client)
			return 200, "ok"
		} else {
			// too many request response
			return 429, "too many requests"
		}
	} else {
		incrementRequestCount(&map[int64]int{}, userId, applicationId, now, client)
		return 200, "ok"

	}
}

func rateInLimit(user *map[int64]int, now time.Time) bool {

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

func incrementRequestCount(redisValue *map[int64]int, userId, applicationId string, now time.Time, client *_redis.Client) {
	key := now.Unix() / 60
	if val, ok := (*redisValue)[key]; ok {
		(*redisValue)[key] = val + 1

	} else {
		(*redisValue)[key] = 1
	}

	b, _ := json.Marshal(redisValue)

	client.Set(getRateLimitingKey(userId, applicationId), string(b), time.Duration(1*time.Hour)).Result()
}

func getLockKey(userId, applicationid string) string {
	return fmt.Sprintf("%s:%s:%s", userId, applicationid, enum.REDIS_LOCK_SUFFIX)
}

func getRateLimitingKey(userId, applicationid string) string {
	return fmt.Sprintf("%s:%s", userId, applicationid)
}
