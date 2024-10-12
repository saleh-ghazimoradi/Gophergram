package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, post *service_modles.Post) error
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

func NewPostRepository(db *sql.DB) Posts {
	return &postRepository{
		db: db,
	}
}
