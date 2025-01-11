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

// CreatePostHandler handles creating a new post.
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		service_models.CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	service_models.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/posts [post]
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

	user := GetUserFromContext(r)

	post := &service_models.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	if err := p.postService.Create(context.Background(), post); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusCreated, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

// GetPostByIdHandler retrieves a specific post by ID.
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	service_models.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/posts/{id} [get]
func (p *PostHandler) GetPostByIdHandler(w http.ResponseWriter, r *http.Request) {
	post := GetPostFromCTX(r)

	comments, err := p.commentService.GetByPostId(context.Background(), post.ID)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err = json.JSONResponse(w, http.StatusOK, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

// UpdatePostHandler updates an existing post.
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Post ID"
//	@Param			payload	body		service_models.UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	service_models.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/posts/{id} [patch]
func (p *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := GetPostFromCTX(r)

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

	if err := json.JSONResponse(w, http.StatusOK, post); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

// DeletePostHandler deletes a post by ID.
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/posts/{id} [delete]
func (p *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := GetPostFromCTX(r)

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

func GetPostFromCTX(r *http.Request) *service_models.Post {
	post, _ := r.Context().Value(PostCtx).(*service_models.Post)
	return post
}

func NewPostHandler(postServer service.PostService, commentService service.CommentService) *PostHandler {
	return &PostHandler{
		postService:    postServer,
		commentService: commentService,
	}
}
