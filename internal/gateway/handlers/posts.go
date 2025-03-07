package handlers

import (
	"context"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"strconv"
)

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

	if err := p.postService.Create(context.Background(), &post); err != nil {
		return helper.InternalServerError(ctx, err)
	}

	return helper.SuccessResponse(ctx, fiber.StatusCreated, "post created successfully", post)
}

func (p *PostHandler) GetPostHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)

	post, err := p.postService.GetById(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	comments, err := p.commentService.GetByPostId(context.Background(), id)
	if err != nil {
		return helper.InternalServerError(ctx, err)
	}

	post.Comments = comments

	return helper.SuccessResponse(ctx, fiber.StatusOK, "success", post)
}

func (p *PostHandler) DeletePostHandler(ctx *fiber.Ctx) error {
	id, _ := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err := p.postService.Delete(context.Background(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			return helper.NotFound(ctx, err)
		default:
			return helper.InternalServerError(ctx, err)
		}
	}

	return helper.SuccessResponse(ctx, fiber.StatusNoContent, "the post successfully deleted", nil)
}

func NewPostHandler(postService service.PostService, commentService service.CommentsService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
	}
}
