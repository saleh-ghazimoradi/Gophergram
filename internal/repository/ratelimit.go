package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RateLimiter interface {
	IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error)
	SetExpiration(ctx context.Context, clientID string, window time.Duration) error
	GetTTL(ctx context.Context, clientID string) (time.Duration, error)
}

type rateLimitRepo struct {
	client *redis.Client
}

func (r *rateLimitRepo) IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error) {
	return r.client.Incr(ctx, clientID).Result()
}

func (r *rateLimitRepo) SetExpiration(ctx context.Context, clientID string, window time.Duration) error {
	return r.client.Expire(ctx, clientID, window).Err()
}

func (r *rateLimitRepo) GetTTL(ctx context.Context, clientID string) (time.Duration, error) {
	return r.client.TTL(ctx, clientID).Result()
}

func NewRateLimitRepo(client *redis.Client) RateLimiter {
	return &rateLimitRepo{client: client}
}
