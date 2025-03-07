package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CommentRepository interface {
	GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error)
	WithTx(tx *sql.Tx) CommentRepository
}

type commentRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (c *commentRepository) GetByPostId(ctx context.Context, postId int64) ([]service_models.Comment, error) {
	comments, err := boiler_models.Comments(
		boiler_models.CommentWhere.PostID.EQ(postId),
		qm.OrderBy("created_at DESC"),
		qm.Load(boiler_models.CommentRels.User),
	).All(ctx, exec(c.dbRead, c.tx))

	if err != nil {
		return nil, err
	}

	serviceComments := make([]service_models.Comment, len(comments))
	for i, comment := range comments {
		serviceComments[i] = service_models.Comment{
			Id:        comment.ID,
			PostId:    comment.PostID,
			UserId:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}

		if comment.R != nil && comment.R.User != nil {
			serviceComments[i].User = service_models.Users{
				ID:       comment.R.User.ID,
				Username: comment.R.User.Username,
			}
		}
	}

	return serviceComments, nil
}

func (c *commentRepository) WithTx(tx *sql.Tx) CommentRepository {
	return &commentRepository{
		dbRead:  c.dbRead,
		dbWrite: c.dbWrite,
		tx:      tx}
}

func NewCommentRepository(dbRead, dbWrite *sql.DB) CommentRepository {
	return &commentRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
