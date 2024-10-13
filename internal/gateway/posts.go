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
	postService service.Posts
}

type postPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (p *Posts) CreatePost(w http.ResponseWriter, r *http.Request) {
	var payload postPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (p *Posts) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := r.Context()

	post, err := p.postService.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func NewPostHandler(postService service.Posts) *Posts {
	return &Posts{
		postService: postService,
	}
}