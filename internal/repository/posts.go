package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

type PostRepository interface {
	Create(ctx context.Context, post *boiler_models.Post) error
	GetById(ctx context.Context, id int64) (*boiler_models.Post, error)
	GetUserFeed(ctx context.Context, userID int64, offset, limit int, search string) ([]*service_models.RawPostWithMetadata, error)
	Update(ctx context.Context, post *boiler_models.Post) error
	Delete(ctx context.Context, id int64) error
	WithTX(tx *sql.Tx) PostRepository
}

type postRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (p *postRepository) Create(ctx context.Context, post *boiler_models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	if err := post.Insert(ctx, exec(p.dbWrite, p.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the post: ", "err", err.Error())
		return err
	}
	return nil
}

func (p *postRepository) GetById(ctx context.Context, id int64) (*boiler_models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	post, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), id)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the post: ", "err", err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}

func (p *postRepository) GetUserFeed(ctx context.Context, userID int64, offset, limit int, search string) ([]*service_models.RawPostWithMetadata, error) {
	const getUserFeedQuery = `
    SELECT 
        p.id, 
        p.user_id, 
        p.title, 
        p.content, 
        p.created_at, 
        p.version, 
        p.tags, 
        u.username, 
        COUNT(c.id) AS comment_count
    FROM posts p
    LEFT JOIN comments c ON c.post_id = p.id
    LEFT JOIN users u ON p.user_id = u.id
    JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
    WHERE (f.user_id = $2 OR p.user_id = $3)
    AND (p.title ILIKE $4 OR p.content ILIKE $4 OR p.tags::text ILIKE $4)
    GROUP BY p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username
    ORDER BY p.created_at DESC
    LIMIT $5 OFFSET $6;
    `

	var rows []*service_models.RawPostWithMetadata

	searchTerm := "%" + search + "%"

	err := queries.Raw(getUserFeedQuery, userID, userID, userID, searchTerm, limit, offset).Bind(ctx, p.dbRead, &rows)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the posts: ", "err", err.Error())
		return nil, err
	}

	sLogger.SLogger.Debug("Fetched rows:", "rows", rows)

	return rows, nil
}

func (p *postRepository) Update(ctx context.Context, post *boiler_models.Post) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	currentPost, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), post.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrsNotFound
		}
		sLogger.SLogger.Error("failed to retrieve the post: ", "err", err.Error())
		return err
	}

	if currentPost.Version != post.Version {
		return errors.New("optimistic lock failed: record has been modified by another transaction")
	}

	post.Version.Int++

	_, err = post.Update(ctx, exec(p.dbWrite, p.tx), boil.Infer())
	if err != nil {
		sLogger.SLogger.Error("failed to update the post: ", "err", err.Error())
		return err
	}
	return nil
}

func (p *postRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	post, err := boiler_models.FindPost(ctx, exec(p.dbRead, p.tx), id)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the post: ", "err", err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrsNotFound
		default:
			return err
		}
	}
	_, err = post.Delete(ctx, exec(p.dbWrite, p.tx))
	if err != nil {
		sLogger.SLogger.Error("failed to delete the post: ", "err", err.Error())
		return err
	}
	return nil
}

func (p *postRepository) WithTX(tx *sql.Tx) PostRepository {
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
