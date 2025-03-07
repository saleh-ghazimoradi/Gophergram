package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type CommentsService interface {
	GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error)
}

type commentsService struct {
	commentsRepository repository.CommentRepository
}

func (c *commentsService) GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error) {
	return c.commentsRepository.GetByPostId(ctx, postId)
}

func NewCommentsService(commentRepository repository.CommentRepository) CommentsService {
	return &commentsService{
		commentsRepository: commentRepository,
	}
}
