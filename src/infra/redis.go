package infra

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries = 3
	retryDelay = 1 * time.Second
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", strings.TrimSpace(os.Getenv("REDIS_HOST")), strings.TrimSpace(os.Getenv("REDIS_PORT"))),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Try to connect with retries
	var err error
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = rdb.Ping(ctx).Result()
		cancel()

		if err == nil {
			log.Println("Successfully connected to Redis")
			return rdb
		}

		log.Printf("Failed to connect to Redis (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	// Even if connection fails, return the client - it will be handled by the repository layer
	log.Printf("Warning: Could not establish Redis connection after %d attempts: %v", maxRetries, err)
	return rdb
}

// IsRedisAvailable checks if Redis is currently available
func IsRedisAvailable(client *redis.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	return err == nil
}
