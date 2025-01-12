package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RateLimitRepository interface {
	IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error)
	SetExpiration(ctx context.Context, clientID string, window time.Duration) error
	GetTTL(ctx context.Context, clientID string) (time.Duration, error)
}

type rateLimitRepository struct {
	client *redis.Client
}

func (r *rateLimitRepository) IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error) {
	return r.client.Incr(ctx, clientID).Result()
}

func (r *rateLimitRepository) SetExpiration(ctx context.Context, clientID string, window time.Duration) error {
	return r.client.Expire(ctx, clientID, window).Err()
}

func (r *rateLimitRepository) GetTTL(ctx context.Context, clientID string) (time.Duration, error) {
	return r.client.TTL(ctx, clientID).Result()
}

func NewRateLimitRepository(client *redis.Client) RateLimitRepository {
	return &rateLimitRepository{
		client: client,
	}
}
