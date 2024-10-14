package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Posts interface {
	Create(ctx context.Context, tx *sql.Tx, post *service_modles.Post) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*service_modles.Post, error)
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
	Update(ctx context.Context, tx *sql.Tx, post *service_modles.Post) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type postRepository struct {
	db *sql.DB
}

func (p *postRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, nil)
}

func (p *postRepository) Create(ctx context.Context, tx *sql.Tx, post *service_modles.Post) error {
	query := `INSERT INTO posts(content, title, user_id, tags) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*service_modles.Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at,  updated_at, tags, version
		FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	var post service_modles.Post
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)
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

func (p *postRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (p *postRepository) Update(ctx context.Context, tx *sql.Tx, post *service_modles.Post) error {
	query := `UPDATE posts SET title = $1, content = $2, version = version + 1 WHERE id = $3 AND version = $4 RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func NewPostRepository(db *sql.DB) Posts {
	return &postRepository{
		db: db,
	}
}
