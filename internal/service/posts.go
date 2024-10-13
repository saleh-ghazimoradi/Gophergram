package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
	GetByID(ctx context.Context, id int64) (*service_modles.Post, error)
}

type postService struct {
	postsRepo repository.Posts
}

func (p *postService) Create(ctx context.Context, post *service_modles.Post) error {
	return p.postsRepo.Create(ctx, post)
}

func (p *postService) GetByID(ctx context.Context, id int64) (*service_modles.Post, error) {
	return p.postsRepo.GetByID(ctx, id)
}

func NewPostService(postsRepo repository.Posts) Posts {
	return &postService{postsRepo: postsRepo}
}
