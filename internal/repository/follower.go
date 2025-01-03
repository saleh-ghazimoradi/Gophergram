package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/config"
)

type FollowerRepository interface {
	Follow(ctx context.Context, followerId, userId int64) error
	Unfollow(ctx context.Context, followerId, userId int64) error
	WithTx(tx *sql.Tx) FollowerRepository
}

type followerRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (f *followerRepository) Follow(ctx context.Context, followerId, userId int64) error {
	query := `INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()
	_, err := f.dbWrite.ExecContext(ctx, query, followerId, userId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrsConflict
		}
	}
	return err
}

func (f *followerRepository) Unfollow(ctx context.Context, followerId, userId int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()
	_, err := f.dbWrite.ExecContext(ctx, query, followerId, userId)

	return err
}

func (f *followerRepository) WithTx(tx *sql.Tx) FollowerRepository {
	return &followerRepository{
		dbWrite: f.dbWrite,
		dbRead:  f.dbRead,
		tx:      tx,
	}
}

func NewFollowerRepository(dbWrite *sql.DB, dbRead *sql.DB) FollowerRepository {
	return &followerRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
