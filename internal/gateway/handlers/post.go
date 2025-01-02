package handlers

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"net/http"
)

type PostKey string

const PostCtx PostKey = "post"

type PostHandler struct {
	postService    service.PostService
	commentService service.CommentService
}

func (p *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_models.CreatePostPayload
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := helper.Validate.Struct(payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	post := &service_models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	if err := p.postService.Create(context.Background(), post); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusCreated, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func (p *PostHandler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCTX(r)

	comments, err := p.commentService.GetByPostId(context.Background(), post.ID)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err = json.WriteJSON(w, http.StatusOK, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func (p *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCTX(r)

	var payload service_models.UpdatePostPayload
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err := helper.Validate.Struct(payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	
	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := p.postService.Update(context.Background(), post); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := json.WriteJSON(w, http.StatusOK, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func (p *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCTX(r)

	if err := p.postService.Delete(context.Background(), post.ID); err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			helper.NotFoundResponse(w, r, err)
		default:
			helper.InternalServerError(w, r, err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func getPostFromCTX(r *http.Request) *service_models.Post {
	post, _ := r.Context().Value(PostCtx).(*service_models.Post)
	return post
}

func NewPostHandler(postServer service.PostService, commentService service.CommentService) *PostHandler {
	return &PostHandler{
		postService:    postServer,
		commentService: commentService,
	}
}
