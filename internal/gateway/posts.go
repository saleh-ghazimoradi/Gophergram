package gateway

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type Posts struct {
	postService    service.Posts
	commentService service.Comments
}

type postPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type updatePayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (p *Posts) CreatePost(w http.ResponseWriter, r *http.Request) {
	var payload postPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	post := &service_modles.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	if err := p.postService.Create(ctx, post); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (p *Posts) GetPost(w http.ResponseWriter, r *http.Request) {
	post := p.GetPostFromCTX(r)

	comments, err := p.commentService.GetByPostID(r.Context(), post.ID)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	post.Comment = comments

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (p *Posts) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	if err := p.postService.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			notFoundResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *Posts) UpdatePost(w http.ResponseWriter, r *http.Request) {
	post := p.GetPostFromCTX(r)

	var payload updatePayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := p.postService.Update(r.Context(), post); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (p *Posts) PostsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := p.postService.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				notFoundResponse(w, r, err)
			default:
				internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (p *Posts) GetPostFromCTX(r *http.Request) *service_modles.Post {
	post, _ := r.Context().Value(postCtx).(*service_modles.Post)
	return post
}

func NewPostHandler(postService service.Posts, commentService service.Comments) Posts {
	return Posts{
		postService:    postService,
		commentService: commentService,
	}
}
