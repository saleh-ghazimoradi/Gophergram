package repository

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PostRepository interface {
	Create(ctx context.Context, post *boiler_models.Post) error
	GetById(ctx context.Context, id int64) (*boiler_models.Post, error)
	Update(ctx context.Context, post *boiler_models.Post) error
	Delete(ctx context.Context, id int64) error
	WithTX(tx *sql.Tx) PostRepository
}

type postRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (p *postRepository) Create(ctx context.Context, post *boiler_models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	if err := post.Insert(ctx, exec(p.dbWrite, p.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the post: ", err)
		return err
	}
	return nil
}

func (p *postRepository) GetById(ctx context.Context, id int64) (*boiler_models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	post, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), id)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the post: ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}

func (p *postRepository) Update(ctx context.Context, post *boiler_models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	currentPost, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), post.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrsNotFound
		}
		sLogger.SLogger.Error("failed to retrieve the post: ", err)
		return err
	}

	if currentPost.Version != post.Version {
		return errors.New("optimistic lock failed: record has been modified by another transaction")
	}

	post.Version.Int++

	_, err = post.Update(ctx, exec(p.dbWrite, p.tx), boil.Infer())
	if err != nil {
		sLogger.SLogger.Error("failed to update the post: ", err)
		return err
	}
	return nil
}

func (p *postRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	post, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), id)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the post: ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrsNotFound
		default:
			return err
		}
	}
	_, err = post.Delete(ctx, exec(p.dbWrite, p.tx))
	if err != nil {
		sLogger.SLogger.Error("failed to delete the post: ", err)
		return err
	}
	return nil
}

func (p *postRepository) WithTX(tx *sql.Tx) PostRepository {
	return &postRepository{
		dbRead:  p.dbRead,
		dbWrite: p.dbWrite,
		tx:      tx,
	}
}

func NewPostRepository(dbRead, dbWrite *sql.DB) PostRepository {
	return &postRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
