package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type PostRepository interface {
	Create(ctx context.Context, post *service_models.Post) error
	GetById(ctx context.Context, id int64) (*service_models.Post, error)
	GetUserFeed(ctx context.Context, id int64, fq service_models.PaginatedFeedQuery) ([]service_models.PostFeed, error)
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

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}

	if err := p.dbWrite.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetById(ctx context.Context, id int64) (*service_models.Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at , version FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	var post service_models.Post

	err := p.dbRead.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version)

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

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

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
	query := `UPDATE posts SET title = $1, content = $2, version = version + 1 WHERE id = $3 AND version = $4 RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	err := p.dbWrite.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrsNotFound
		default:
			return err
		}
	}
	return nil
}

func (p *postRepository) GetUserFeed(ctx context.Context, id int64, fq service_models.PaginatedFeedQuery) ([]service_models.PostFeed, error) {
	query := `
	   SELECT
	       p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
	       u.username,
	       COUNT(c.id) AS comments_count
	   FROM posts p
	   LEFT JOIN comments c ON c.post_id = p.id
	   LEFT JOIN users u ON p.user_id = u.id
	   JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
	   WHERE f.user_id = $1 OR p.user_id = $1
	   GROUP BY p.id, u.username
	   ORDER BY p.created_at ` + fq.Sort + `
	   LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	rows, err := p.dbRead.QueryContext(ctx, query, id, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []service_models.PostFeed
	for rows.Next() {
		var ps service_models.PostFeed
		err := rows.Scan(
			&ps.ID,
			&ps.UserID,
			&ps.Title,
			&ps.Content,
			&ps.CreatedAt,
			&ps.Version,
			pq.Array(&ps.Tags),
			&ps.User.Username,
			&ps.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, ps)
	}
	return feed, nil
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
