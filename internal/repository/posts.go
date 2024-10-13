package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
	GetByID(ctx context.Context, id int64) (*service_modles.Post, error)
}

type postRepository struct {
	db *sql.DB
}

func (p *postRepository) Create(ctx context.Context, post *service_modles.Post) error {
	query := `INSERT INTO posts(content, title, user_id, tags) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := p.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetByID(ctx context.Context, id int64) (*service_modles.Post, error) {
	query := `SELECT id, content, title, user_id, created_at, updated_at,tags FROM posts WHERE id = $1`
	var post service_modles.Post
	err := p.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, &post.CreatedAt, &post.UpdatedAt, pq.Array(&post.Tags))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func NewPostRepository(db *sql.DB) Posts {
	return &postRepository{
		db: db,
	}
}
