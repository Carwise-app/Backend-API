package infra

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", strings.TrimSpace(os.Getenv("REDIS_HOST")), strings.TrimSpace(os.Getenv("REDIS_PORT"))),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("failed to connect to Redis: %v", err)
	}
	return rdb
}
