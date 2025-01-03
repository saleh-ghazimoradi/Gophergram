package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type CommentService interface {
	GetByPostId(ctx context.Context, id int64) ([]service_models.Comment, error)
	Create(ctx context.Context, comment *service_models.Comment) error
}

type commentService struct {
	commentRepo repository.CommentRepository
}

func (c *commentService) GetByPostId(ctx context.Context, id int64) ([]service_models.Comment, error) {
	return c.commentRepo.GetByPostId(ctx, id)
}

func (c *commentService) Create(ctx context.Context, comment *service_models.Comment) error {
	return c.commentRepo.Create(ctx, comment)
}

func NewCommentService(commentRepo repository.CommentRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
	}
}
