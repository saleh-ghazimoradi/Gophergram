package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"time"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type RateLimiter interface {
	IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error)
	SetExpiration(ctx context.Context, clientID string, window time.Duration) error
	GetTTL(ctx context.Context, clientID string) (time.Duration, error)
	IsAllowed(ctx context.Context, clientID string, limit int, window time.Duration) (bool, time.Duration, error)
}

type rateLimitService struct {
	rateLimitRepo repository.RateLimiter
}

func (r *rateLimitService) IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error) {
	return r.rateLimitRepo.IncrementRequestCount(ctx, clientID, window)
}

func (r *rateLimitService) SetExpiration(ctx context.Context, clientID string, window time.Duration) error {
	return r.rateLimitRepo.SetExpiration(ctx, clientID, window)
}

func (r *rateLimitService) GetTTL(ctx context.Context, clientID string) (time.Duration, error) {
	return r.rateLimitRepo.GetTTL(ctx, clientID)
}

func (r *rateLimitService) IsAllowed(ctx context.Context, clientID string, limit int, window time.Duration) (bool, time.Duration, error) {
	count, err := r.rateLimitRepo.IncrementRequestCount(ctx, clientID, window)
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		if err := r.rateLimitRepo.SetExpiration(ctx, clientID, window); err != nil {
			return false, 0, err
		}
	}

	if count > int64(limit) {
		ttl, _ := r.rateLimitRepo.GetTTL(ctx, clientID)
		return false, ttl, ErrRateLimitExceeded
	}

	return true, 0, nil
}

func NewRateLimitService(rateLimitRepo repository.RateLimiter) RateLimiter {
	return &rateLimitService{rateLimitRepo: rateLimitRepo}
}
