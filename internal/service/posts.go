package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostService interface {
	Create(ctx context.Context, post *dto.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
}

type postService struct {
	postRepository repository.PostRepository
}

func (p *postService) Create(ctx context.Context, input *dto.Post) error {
	return p.postRepository.Create(ctx, &service_models.Post{
		Content: input.Content,
		Title:   input.Title,
		Tags:    input.Tags,
		UserID:  1,
	})
}

func (p *postService) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
	return p.postRepository.GetById(ctx, id)
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepository,
	}
}
