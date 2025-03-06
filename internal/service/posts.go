package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostService interface {
	Create(ctx context.Context, post *dto.Post) error
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

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepository,
	}
}
