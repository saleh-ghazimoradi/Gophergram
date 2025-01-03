package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type CommentRepository interface {
	GetByPostId(ctx context.Context, id int64) ([]service_models.Comment, error)
	Create(ctx context.Context, comment *service_models.Comment) error
	WithTx(tx *sql.Tx) CommentRepository
}

type commentRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (c *commentRepository) GetByPostId(ctx context.Context, id int64) ([]service_models.Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c JOIN users on users.id = c.user_id WHERE c.post_id = $1 ORDER BY c.created_at DESC;`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	rows, err := c.dbRead.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]service_models.Comment, 0)
	for rows.Next() {
		var comment service_models.Comment
		comment.User = service_models.User{}
		err = rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.UserID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *commentRepository) Create(ctx context.Context, comment *service_models.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()
	err := c.dbWrite.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *commentRepository) WithTx(tx *sql.Tx) CommentRepository {
	return &commentRepository{
		dbRead:  c.dbRead,
		dbWrite: c.dbWrite,
		tx:      tx,
	}
}

func NewCommentRepository(dbRead, dbWrite *sql.DB) CommentRepository {
	return &commentRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
