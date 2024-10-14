package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"log"
)

type Comments interface {
	Create(ctx context.Context, comments *service_modles.Comments) error
	GetByPostID(ctx context.Context, id int64) ([]service_modles.Comments, error)
}

type commentService struct {
	commentRepo repository.Comments
}

func (c *commentService) GetByPostID(ctx context.Context, id int64) ([]service_modles.Comments, error) {
	tx, err := c.commentRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
			}
		}
	}()
	comment, err := c.commentRepo.GetByPostID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return nil, commitErr
	}
	return comment, nil
}

func (c *commentService) Create(ctx context.Context, comments *service_modles.Comments) error {
	tx, err := c.commentRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err := c.commentRepo.Create(ctx, tx, comments); err != nil {
		return err
	}
	return nil
}

func NewCommentService(commentRepo repository.Comments) Comments {
	return &commentService{commentRepo: commentRepo}
}
