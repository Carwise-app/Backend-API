package infra

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository() *TokenRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", strings.TrimSpace(os.Getenv("REDIS_HOST")), strings.TrimSpace(os.Getenv("REDIS_PORT"))),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
	return &TokenRepository{client: rdb}
}

func (r *TokenRepository) IsTokenBlackListed(token string) (bool, error) {
	val, err := r.client.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token in Redis: %v", err)
	}

	if val == "blacklisted" {
		return true, nil
	}

	return false, nil
}

func (r *TokenRepository) AddTokenBlackList(token string) error {
	err := r.client.Set(context.Background(), token, "blacklisted", 0).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist in Redis: %v", err)
	}

	return nil
}
