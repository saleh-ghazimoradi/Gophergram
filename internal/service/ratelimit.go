package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"time"
)

type RateLimitService interface {
	IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error)
	SetExpiration(ctx context.Context, clientID string, window time.Duration) error
	GetTTL(ctx context.Context, clientID string) (time.Duration, error)
	IsAllowed(ctx context.Context, clientID string, limit int, window time.Duration) (bool, time.Duration, error)
}

type rateLimitService struct {
	rateLimitRepository repository.RateLimitRepository
}

func (r *rateLimitService) IncrementRequestCount(ctx context.Context, clientID string, window time.Duration) (int64, error) {
	return r.rateLimitRepository.IncrementRequestCount(ctx, clientID, window)
}

func (r *rateLimitService) SetExpiration(ctx context.Context, clientID string, window time.Duration) error {
	return r.rateLimitRepository.SetExpiration(ctx, clientID, window)
}

func (r *rateLimitService) GetTTL(ctx context.Context, clientID string) (time.Duration, error) {
	return r.rateLimitRepository.GetTTL(ctx, clientID)
}

func (r *rateLimitService) IsAllowed(ctx context.Context, clientID string, limit int, window time.Duration) (bool, time.Duration, error) {
	count, err := r.rateLimitRepository.IncrementRequestCount(ctx, clientID, window)
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		if err := r.rateLimitRepository.SetExpiration(ctx, clientID, window); err != nil {
			return false, 0, err
		}
	}

	if count > int64(limit) {
		ttl, _ := r.rateLimitRepository.GetTTL(ctx, clientID)
		return false, ttl, repository.ErrRateLimitExceeded
	}

	return true, 0, nil
}

func NewRateLimitService(rateLimitRepository repository.RateLimitRepository) RateLimitService {
	return &rateLimitService{
		rateLimitRepository: rateLimitRepository,
	}
}
