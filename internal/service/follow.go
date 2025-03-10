package service

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"time"
)

type FollowService interface {
	Follow(ctx context.Context, followerId, userId int64) error
	Unfollow(ctx context.Context, followerId, userId int64) error
}

type followService struct {
	followerRepository repository.FollowerRepository
}

func (f *followService) Follow(ctx context.Context, followerId, userId int64) error {
	if followerId == userId {
		return errors.New("you cannot follow yourself")
	}

	follower := &boiler_models.Follower{
		FollowerID: followerId,
		UserID:     userId,
		CreatedAt:  time.Now(),
	}

	if err := f.followerRepository.Follow(ctx, follower); err != nil {
		return err
	}

	return nil
}

func (f *followService) Unfollow(ctx context.Context, followerId, userId int64) error {
	if followerId == userId {
		return errors.New("you cannot unfollow yourself")
	}

	follower := &boiler_models.Follower{
		FollowerID: followerId,
		UserID:     userId,
		CreatedAt:  time.Now(),
	}

	if err := f.followerRepository.Unfollow(ctx, follower); err != nil {
		return err
	}

	return nil
}

func NewFollowService(followerRepository repository.FollowerRepository) FollowService {
	return &followService{
		followerRepository: followerRepository,
	}
}
