package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
	GetByID(ctx context.Context, id int64) (*service_modles.Post, error)
	Delete(ctx context.Context, id int64) error
}

type postService struct {
	postRepo repository.Posts
}

func (p *postService) Create(ctx context.Context, post *service_modles.Post) error {
	return p.postRepo.Create(ctx, post)
}

func (p *postService) GetByID(ctx context.Context, id int64) (*service_modles.Post, error) {
	return p.postRepo.GetByID(ctx, id)
}

func (p *postService) Delete(ctx context.Context, id int64) error {
	return p.postRepo.Delete(ctx, id)
}

func NewPostService(postsRepo repository.Posts) Posts {
	return &postService{postRepo: postsRepo}
}
