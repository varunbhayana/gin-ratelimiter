package db

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	_redis "github.com/go-redis/redis"
	_ "github.com/lib/pq" //import postgres
	"github.com/varunbhayana/rate-limiting/conf"
)

//DB ...
type DB struct {
	*sql.DB
}

var db *gorp.DbMap

//Init ...
func Init() {

}

//RedisClient ...
var RedisClient *_redis.Client

//InitRedis ...
func InitRedis(params ...string) {

	var redisHost = conf.D.REDIS_HOST
	var redisPassword = conf.D.REDIS_PASSWORD
	//db, _ := strconv.Atoi(params[0])

	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
	})
}

//GetRedis ...
func GetRedis() *_redis.Client {
	return RedisClient
}
