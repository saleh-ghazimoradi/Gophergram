package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Comments interface {
	Create(ctx context.Context, comments *service_modles.Comments) error
	GetByPostID(ctx context.Context, id int64) ([]service_modles.Comments, error)
}

type commentService struct {
	commentRepo repository.Comments
	db          *sql.DB
}

func (c *commentService) GetByPostID(ctx context.Context, id int64) ([]service_modles.Comments, error) {
	return c.commentRepo.GetByPostID(ctx, id)
}

func (c *commentService) Create(ctx context.Context, comments *service_modles.Comments) error {
	_, err := withTransaction(ctx, c.db, func(tx *sql.Tx) (struct{}, error) {
		if err := c.commentRepo.Create(ctx, tx, comments); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, nil
	})
	return err
}

func NewCommentService(commentRepo repository.Comments) Comments {
	return &commentService{commentRepo: commentRepo}
}
