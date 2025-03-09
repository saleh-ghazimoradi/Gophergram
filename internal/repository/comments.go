package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CommentRepository interface {
	GetByPostId(ctx context.Context, postId int64) ([]*boiler_models.Comment, error)
	Create(ctx context.Context, comment *boiler_models.Comment) error
	WithTx(tx *sql.Tx) CommentRepository
}

type commentRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (c *commentRepository) GetByPostId(ctx context.Context, postId int64) ([]*boiler_models.Comment, error) {
	comments, err := boiler_models.Comments(
		boiler_models.CommentWhere.PostID.EQ(postId),
		qm.OrderBy("created_at DESC"),
		qm.Load(boiler_models.CommentRels.User),
	).All(ctx, exec(c.dbRead, c.tx))

	if err != nil {
		sLogger.SLogger.Error("failed to fetch comments by post ID", err)
		return nil, err
	}

	return comments, nil
}

func (c *commentRepository) Create(ctx context.Context, comment *boiler_models.Comment) error {
	if err := comment.Insert(ctx, c.dbWrite, boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the comment", err)
		return err
	}
	return nil
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
