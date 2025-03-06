package repository

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type PostRepository interface {
	Create(ctx context.Context, post *service_models.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
	WithTX(tx *sql.Tx) PostRepository
}

type postRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (p *postRepository) Create(ctx context.Context, post *service_models.Post) error {
	newPost := &boiler_models.Post{
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
	}

	if err := newPost.Insert(ctx, exec(p.dbWrite, p.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the post: ", err)
		return err
	}
	
	return nil
}

func (p *postRepository) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
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

	return &service_models.Post{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
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
