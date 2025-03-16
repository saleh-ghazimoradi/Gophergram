package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type CommentsService interface {
	GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error)
	Create(ctx context.Context, input *dto.Comment) error
}

type commentsService struct {
	commentsRepository repository.CommentRepository
}

func (c *commentsService) GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error) {
	boilerComments, err := c.commentsRepository.GetByPostId(ctx, postId)
	if err != nil {
		return nil, err
	}

	serviceComments := make([]service_models.Comment, len(boilerComments))
	for i, comment := range boilerComments {
		serviceComments[i] = service_models.Comment{
			Id:        comment.ID,
			PostId:    comment.PostID,
			UserId:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}

		if comment.R != nil && comment.R.User != nil {
			serviceComments[i].User = service_models.Users{
				ID:       comment.R.User.ID,
				Username: comment.R.User.Username,
			}
		}
	}

	return serviceComments, nil
}

func (c *commentsService) Create(ctx context.Context, input *dto.Comment) error {
	newComment := &boiler_models.Comment{
		Content: input.Content,
	}
	return c.commentsRepository.Create(ctx, newComment)
}

func NewCommentsService(commentRepository repository.CommentRepository) CommentsService {
	return &commentsService{
		commentsRepository: commentRepository,
	}
}
