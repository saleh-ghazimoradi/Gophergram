package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostRepository interface {
	Create(ctx context.Context, post *service_models.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *service_models.Post) error
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

func (p *postRepository) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`

	var post service_models.Post

	err := p.dbRead.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (p *postRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	result, err := p.dbWrite.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrsNotFound
	}
	return nil
}

func (p *postRepository) Update(ctx context.Context, post *service_models.Post) error {
	query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`

	_, err := p.dbWrite.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
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
