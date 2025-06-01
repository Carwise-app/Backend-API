package infra

import (
	"carwise"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	BRANDS_CACHE_KEY   = "brands:with_details"
	BRANDS_CACHE_TTL   = 5 * time.Minute
	LISTINGS_CACHE_TTL = 5 * time.Minute
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{client: ConnectRedis()}
}

func (r *RedisRepository) IsTokenBlackListed(token string) (bool, error) {
	if !IsRedisAvailable(r.client) {
		return false, fmt.Errorf("redis is not available")
	}

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

func (r *RedisRepository) AddTokenBlackList(token string) error {
	if !IsRedisAvailable(r.client) {
		return fmt.Errorf("redis is not available")
	}

	err := r.client.Set(context.Background(), token, "blacklisted", 0).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist in Redis: %v", err)
	}

	return nil
}

func (r *RedisRepository) SetBrandsWithDetails(brands []carwise.BrandWithDetails) error {
	if !IsRedisAvailable(r.client) {
		return fmt.Errorf("redis is not available")
	}

	data, err := json.Marshal(brands)
	if err != nil {
		return fmt.Errorf("failed to marshal brands: %v", err)
	}

	err = r.client.Set(context.Background(), BRANDS_CACHE_KEY, data, BRANDS_CACHE_TTL).Err()
	if err != nil {
		return fmt.Errorf("failed to cache brands in Redis: %v", err)
	}

	return nil
}

func (r *RedisRepository) GetBrandsWithDetails() ([]carwise.BrandWithDetails, error) {
	if !IsRedisAvailable(r.client) {
		return nil, fmt.Errorf("redis is not available")
	}

	data, err := r.client.Get(context.Background(), BRANDS_CACHE_KEY).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get brands from Redis: %v", err)
	}

	var brands []carwise.BrandWithDetails
	if err := json.Unmarshal(data, &brands); err != nil {
		return nil, fmt.Errorf("failed to unmarshal brands: %v", err)
	}

	return brands, nil
}
