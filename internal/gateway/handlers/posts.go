package handlers

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type PostHandler struct {
	postService    service.PostService
	commentService service.CommentsService
}

func (p *PostHandler) CreatePostHandler(ctx *fiber.Ctx) error {
	var post dto.Post
	if err := ctx.BodyParser(&post); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := helper.Validator.Struct(post); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := p.postService.Create(ctx.Context(), &post); err != nil {
		return helper.InternalServerError(ctx, err)
	}

	return helper.SuccessResponse(ctx, fiber.StatusCreated, "post created successfully", post)
}

func (p *PostHandler) GetPostHandler(ctx *fiber.Ctx) error {
	post := p.GetPostFromContext(ctx)

	comments, err := p.commentService.GetByPostId(context.Background(), post.ID)
	if err != nil {
		return helper.InternalServerError(ctx, err)
	}

	post.Comments = comments

	return helper.SuccessResponse(ctx, fiber.StatusOK, "success", post)
}

func (p *PostHandler) UpdatePostHandler(ctx *fiber.Ctx) error {
	post := p.GetPostFromContext(ctx)
	var payload dto.UpdatePost
	if err := ctx.BodyParser(&payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := helper.Validator.Struct(payload); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := p.postService.Update(ctx.Context(), &payload, post); err != nil {
		switch {
		case errors.Is(err, errors.New("conflict: record has been modified by another transaction")):
			return helper.ConflictResponse(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusOK, "success", post)
}

func (p *PostHandler) DeletePostHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err := p.postService.Delete(ctx.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusNoContent, "the post successfully deleted", nil)
}

func (p *PostHandler) PostsContextMiddleware(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return helper.BadRequest(ctx, err)
	}

	post, err := p.postService.GetById(ctx.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	ctx.Locals(postCtx, post)
	return ctx.Next()
}

func (p *PostHandler) GetPostFromContext(ctx *fiber.Ctx) *service_models.Post {
	post, ok := ctx.Locals(postCtx).(*service_models.Post)
	if !ok {
		sLogger.SLogger.Warn("Post not found in context")
		return nil
	}
	return post
}

func NewPostHandler(postService service.PostService, commentService service.CommentsService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
	}
}
