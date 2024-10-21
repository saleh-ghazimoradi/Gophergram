package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
	GetByID(ctx context.Context, id int64) (*service_modles.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *service_modles.Post) error
	GetUserFeed(ctx context.Context, id int64, fq service_modles.PaginatedFeedQuery) ([]service_modles.PostWithMetaData, error)
}

type postService struct {
	postRepo    repository.Posts
	commentRepo repository.Comments
	db          *sql.DB
}

func (p *postService) Create(ctx context.Context, post *service_modles.Post) error {
	_, err := withTransaction(ctx, p.db, func(tx *sql.Tx) (struct{}, error) {
		return struct{}{}, p.postRepo.Create(ctx, tx, post)
	})
	return err
}

func (p *postService) GetByID(ctx context.Context, id int64) (*service_modles.Post, error) {
	return withTransaction(ctx, p.db, func(tx *sql.Tx) (*service_modles.Post, error) {
		return p.postRepo.GetByID(ctx, tx, id)
	})
}

func (p *postService) Delete(ctx context.Context, id int64) error {
	_, err := withTransaction(ctx, p.db, func(tx *sql.Tx) (struct{}, error) {
		return struct{}{}, p.postRepo.Delete(ctx, tx, id)
	})
	return err
}

func (p *postService) Update(ctx context.Context, post *service_modles.Post) error {
	_, err := withTransaction(ctx, p.db, func(tx *sql.Tx) (struct{}, error) {
		return struct{}{}, p.postRepo.Update(ctx, tx, post)
	})
	return err
}

func (p *postService) GetUserFeed(ctx context.Context, id int64, fq service_modles.PaginatedFeedQuery) ([]service_modles.PostWithMetaData, error) {
	return p.postRepo.GetUserFeed(ctx, id, fq)
}

func NewPostService(postsRepo repository.Posts, commentRepo repository.Comments, db *sql.DB) Posts {
	return &postService{postRepo: postsRepo, commentRepo: commentRepo, db: db}
}
