package repository

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

type FollowerRepository interface {
	Follow(ctx context.Context, follower *boiler_models.Follower) error
	Unfollow(ctx context.Context, follower *boiler_models.Follower) error
	WithTx(tx *sql.Tx) FollowerRepository
}

type followerRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (f *followerRepository) Follow(ctx context.Context, follower *boiler_models.Follower) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := follower.Insert(ctx, exec(f.dbWrite, f.tx), boil.Infer())
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

func (f *followerRepository) Unfollow(ctx context.Context, follower *boiler_models.Follower) error {
	rowsAffected, err := follower.Delete(ctx, exec(f.dbWrite, f.tx))
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected; follower relationship may not exist")
	}

	return nil
}

func (f *followerRepository) WithTx(tx *sql.Tx) FollowerRepository {
	return &followerRepository{
		dbRead:  f.dbRead,
		dbWrite: f.dbWrite,
		tx:      tx,
	}
}

func NewFollowerRepository(dbRead *sql.DB, dbWrite *sql.DB) FollowerRepository {
	return &followerRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
