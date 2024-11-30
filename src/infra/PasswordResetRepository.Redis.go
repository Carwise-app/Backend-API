package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PasswordResetRepository struct {
	client *redis.Client
}

func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{client: ConnectRedis()}
}

func (r *PasswordResetRepository) SaveResetCode(email, code string, ttl time.Duration) error {
	key := fmt.Sprintf("password-reset:%s", email)
	err := r.client.Set(context.Background(), key, code, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save password reset code to Redis: %v", err)
	}

	return nil
}

func (r *PasswordResetRepository) VerifyResetCode(email, code string) (bool, error) {
	key := fmt.Sprintf("password-reset:%s", email)
	storedCode, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to retrieve password reset code from Redis: %v", err)
	}

	return storedCode == code, nil
}

func (r *PasswordResetRepository) DeleteResetCode(email string) error {
	key := fmt.Sprintf("password-reset:%s", email)
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete password reset code from Redis: %v", err)
	}

	return nil
}
