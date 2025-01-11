package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type CacheService interface {
	Get(ctx context.Context, id int64) (*service_models.User, error)
	Set(ctx context.Context, user *service_models.User) error
}

type cacheService struct {
	cacheRepository repository.CacheRepository
}

func (s *cacheService) Get(ctx context.Context, id int64) (*service_models.User, error) {
	return s.cacheRepository.Get(ctx, id)
}

func (s *cacheService) Set(ctx context.Context, user *service_models.User) error {
	return s.cacheRepository.Set(ctx, user)
}

func NewCacheService(cacheRepository repository.CacheRepository) CacheService {
	return &cacheService{
		cacheRepository: cacheRepository,
	}
}
