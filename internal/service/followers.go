package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
)

type Follow interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

type followService struct {
	followRepo repository.Follow
	db         *sql.DB
}

func (s *followService) Follow(ctx context.Context, followerID, userID int64) error {
	_, err := withTransaction(ctx, s.db, func(tx *sql.Tx) (struct{}, error) {
		if err := s.followRepo.Follow(ctx, tx, followerID, userID); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	})
	return err
}

func (s *followService) Unfollow(ctx context.Context, followerID, userID int64) error {
	_, err := withTransaction(ctx, s.db, func(tx *sql.Tx) (struct{}, error) {
		if err := s.followRepo.Unfollow(ctx, tx, followerID, userID); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	})
	return err
}

func NewFollowService(followRepo repository.Follow, db *sql.DB) Follow {
	return &followService{
		followRepo: followRepo,
		db:         db,
	}
}
