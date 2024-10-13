package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Comments interface {
	GetByPostID(ctx context.Context, postID int64) ([]service_modles.Comments, error)
}

type commentService struct {
	commentRepo repository.Comments
}

func (c *commentService) GetByPostID(ctx context.Context, postID int64) ([]service_modles.Comments, error) {
	return c.commentRepo.GetByPostID(ctx, postID)
}

func NewCommentService(commentRepo repository.Comments) Comments {
	return &commentService{commentRepo: commentRepo}
}
