package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
)

type FollowerService interface {
	Follow(ctx context.Context, followerId, userId int64) error
	Unfollow(ctx context.Context, followerId, userId int64) error
}

type followerService struct {
	followerRepo repository.FollowerRepository
}

func (s *followerService) Follow(ctx context.Context, followerId, userId int64) error {
	return s.followerRepo.Follow(ctx, followerId, userId)
}

func (s *followerService) Unfollow(ctx context.Context, followerId, userId int64) error {
	return s.followerRepo.Unfollow(ctx, followerId, userId)
}

func NewFollowerService(followerRepo repository.FollowerRepository) FollowerService {
	return &followerService{
		followerRepo: followerRepo,
	}
}
