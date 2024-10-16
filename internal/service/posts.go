package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"log"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
	GetByID(ctx context.Context, id int64) (*service_modles.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *service_modles.Post) error
	GetUserFeed(ctx context.Context, id int64, fq service_modles.PaginatedFeedQuery) ([]service_modles.PostWithMetaData, error)
}

type postService struct {
	postRepo    repository.Posts
	commentRepo repository.Comments
}

func (p *postService) Create(ctx context.Context, post *service_modles.Post) error {
	tx, err := p.postRepo.BeginTx(ctx)
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
	if err := p.postRepo.Create(ctx, tx, post); err != nil {
		return err
	}
	return nil
}

func (p *postService) GetByID(ctx context.Context, id int64) (*service_modles.Post, error) {
	tx, err := p.postRepo.BeginTx(ctx)
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
	post, err := p.postRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if postErr := tx.Commit(); postErr != nil {
		return nil, postErr
	}

	return post, nil
}

func (p *postService) Delete(ctx context.Context, id int64) error {
	tx, err := p.postRepo.BeginTx(ctx)
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

	if err := p.postRepo.Delete(ctx, tx, id); err != nil {
		return err
	}
	return nil
}

func (p *postService) Update(ctx context.Context, post *service_modles.Post) error {
	tx, err := p.postRepo.BeginTx(ctx)
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
	if err := p.postRepo.Update(ctx, tx, post); err != nil {
		return err
	}
	return nil
}

func (p *postService) GetUserFeed(ctx context.Context, id int64, fq service_modles.PaginatedFeedQuery) ([]service_modles.PostWithMetaData, error) {
	return p.postRepo.GetUserFeed(ctx, id, fq)
}

func NewPostService(postsRepo repository.Posts, commentRepo repository.Comments) Posts {
	return &postService{postRepo: postsRepo, commentRepo: commentRepo}
}
