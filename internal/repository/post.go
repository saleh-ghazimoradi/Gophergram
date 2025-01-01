package repository

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostRepository interface {
	Create(ctx context.Context, post *service_models.Post) error
	WithTx(tx *sql.Tx) PostRepository
}

type postRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (p *postRepository) Create(ctx context.Context, post *service_models.Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}

	if err := p.dbWrite.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (p *postRepository) WithTx(tx *sql.Tx) PostRepository {
	return &postRepository{
		dbRead:  p.dbRead,
		dbWrite: p.dbWrite,
		tx:      tx,
	}
}

func NewPostRepository(dbRead, dbWrite *sql.DB) PostRepository {
	return &postRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
