package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/config"
)

type Follow interface {
	Follow(ctx context.Context, tx *sql.Tx, followerID, userID int64) error
	Unfollow(ctx context.Context, tx *sql.Tx, followerID, userID int64) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type followRepository struct {
	DB *sql.DB
}

func (f *followRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return f.DB.BeginTx(ctx, nil)
}

func (f *followRepository) Follow(ctx context.Context, tx *sql.Tx, followerID, userID int64) error {
	query := `
	INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)
`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return nil
}

func (f *followRepository) Unfollow(ctx context.Context, tx *sql.Tx, followerID, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id= $2`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userID, followerID)
	return err
}

func NewFollowRepository(db *sql.DB) Follow {
	return &followRepository{DB: db}
}
