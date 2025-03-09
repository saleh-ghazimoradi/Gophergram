package service

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"time"
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
	boilerPost := &boiler_models.Post{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    1,
		CreatedAt: time.Now(),
	}

	if err := p.postRepository.Create(ctx, boilerPost); err != nil {
		return err
	}

	return nil
}

func (p *postService) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
	boilerPost, err := p.postRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &service_models.Post{
		ID:        boilerPost.ID,
		Title:     boilerPost.Title,
		Content:   boilerPost.Content,
		UserID:    boilerPost.UserID,
		CreatedAt: boilerPost.CreatedAt,
		UpdatedAt: boilerPost.UpdatedAt,
		Version:   boilerPost.Version.Int,
	}, nil
}

func (p *postService) Update(ctx context.Context, input *dto.UpdatePost, post *service_models.Post) error {
	existingPost, err := p.postRepository.GetById(ctx, post.ID)
	if err != nil {
		return err
	}

	if input.Title != nil {
		existingPost.Title = *input.Title
	}

	if input.Content != nil {
		existingPost.Content = *input.Content
	}

	if err = p.postRepository.Update(ctx, existingPost); err != nil {
		if errors.Is(err, errors.New("optimistic lock failed")) {
			return errors.New("conflict: record has been modified by another transaction")
		}
		return err
	}

	post.Version = existingPost.Version.Int

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
