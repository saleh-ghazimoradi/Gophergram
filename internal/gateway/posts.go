package gateway

import (
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
	"strconv"
)

type Posts struct {
	postService    service.Posts
	commentService service.Comments
}

type postPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
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

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (p *Posts) GetPost(w http.ResponseWriter, r *http.Request) {
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

	comments, err := p.commentService.GetByPostID(ctx, id)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func NewPostHandler(postService service.Posts) *Posts {
	return &Posts{
		postService: postService,
	}
}
