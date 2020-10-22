package db

import (
	"database/sql"
	"os"

	"github.com/go-gorp/gorp"
	_redis "github.com/go-redis/redis"
	_ "github.com/lib/pq" //import postgres
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

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

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
