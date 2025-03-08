package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostService interface {
	Create(ctx context.Context, input *dto.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
	Update(ctx context.Context, input *dto.UpdatePost, post *service_models.Post) error
	Delete(ctx context.Context, id int64) error
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

func (p *postService) Update(ctx context.Context, input *dto.UpdatePost, post *service_models.Post) error {
	// Map fields from DTO to the service model
	if input.Title != nil {
		post.Title = *input.Title
	}

	if input.Content != nil {
		post.Content = *input.Content
	}

	// Update the post in the repository
	err := p.postRepository.Update(ctx, post)
	if err != nil {
		return err
	}

	return nil
}

func (p *postService) Delete(ctx context.Context, id int64) error {
	return p.postRepository.Delete(ctx, id)
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepository,
	}
}
