package enum

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	REDIS_LOCK_SUFFIX = "_lock"
)

func ReadEnv(key string) string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}
	return os.Getenv(key)
}
