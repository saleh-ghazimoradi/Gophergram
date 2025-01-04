package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
)

type PostService interface {
	Create(ctx context.Context, post *service_models.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
	GetUserFeed(ctx context.Context, id int64, fq service_models.PaginatedFeedQuery) ([]service_models.PostFeed, error)
	Update(ctx context.Context, post *service_models.Post) error
	Delete(ctx context.Context, id int64) error
}

type postService struct {
	postRepo repository.PostRepository
	db       *sql.DB
}

func (p *postService) Create(ctx context.Context, post *service_models.Post) error {
	return utils.WithTransaction(ctx, p.db, func(tx *sql.Tx) error {
		userRepoWithTx := p.postRepo.WithTx(tx)
		return userRepoWithTx.Create(ctx, post)
	})
}

func (p *postService) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
	return p.postRepo.GetById(ctx, id)
}

func (p *postService) Update(ctx context.Context, post *service_models.Post) error {
	return p.postRepo.Update(ctx, post)
}

func (p *postService) Delete(ctx context.Context, id int64) error {
	return p.postRepo.Delete(ctx, id)
}

func (p *postService) GetUserFeed(ctx context.Context, id int64, fq service_models.PaginatedFeedQuery) ([]service_models.PostFeed, error) {
	return p.postRepo.GetUserFeed(ctx, id, fq)
}

func NewPostService(postRepo repository.PostRepository, db *sql.DB) PostService {
	return &postService{
		postRepo: postRepo,
		db:       db,
	}
}
