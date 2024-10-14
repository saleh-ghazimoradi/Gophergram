package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"log"
)

type Follow interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

type followService struct {
	followRepo repository.Follow
}

func (s *followService) Follow(ctx context.Context, followerID, userID int64) error {
	tx, err := s.followRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
			}
		}
	}()

	err = s.followRepo.Follow(ctx, tx, followerID, userID)
	if err != nil {
		return err
	}
	if followerErr := tx.Commit(); followerErr != nil {
		return followerErr
	}
	return nil
}

func (s *followService) Unfollow(ctx context.Context, followerID, userID int64) error {
	tx, err := s.followRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
			}
		}
	}()
	err = s.followRepo.Unfollow(ctx, tx, followerID, userID)
	if err != nil {
		return err
	}
	if unfollowerErr := tx.Commit(); unfollowerErr != nil {
		return unfollowerErr
	}
	return nil
}

func NewFollowService(followRepo repository.Follow) Follow {
	return &followService{
		followRepo: followRepo,
	}
}
