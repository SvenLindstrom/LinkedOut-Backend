package databases

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func Redis_init() *redis.Client {

	user := os.Getenv("REDIS_USER")
	pwd := os.Getenv("REDIS_PASSWORD")
	addr := os.Getenv("REDIS_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Username: user,
		Password: pwd,
		Addr:     addr,
	})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to Redis DB at %s", addr)
	}

	return rdb
}
