package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Comments interface {
	GetByPostID(ctx context.Context, postID int64) ([]service_modles.Comments, error)
}

type commentRepository struct {
	DB *sql.DB
}

func (c *commentRepository) GetByPostID(ctx context.Context, postID int64) ([]service_modles.Comments, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments c JOIN users on users.id = c.user_id WHERE c.post_id = $1 ORDER BY c.created_at DESC `

	rows, err := c.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]service_modles.Comments, 0)

	for rows.Next() {
		var c service_modles.Comments
		c.User = service_modles.Users{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func NewCommentRepository(db *sql.DB) Comments {
	return &commentRepository{DB: db}
}
