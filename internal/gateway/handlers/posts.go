package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func (p *PostHandler) CreatePostHandler(ctx *fiber.Ctx) error {
	var post dto.Post
	if err := ctx.BodyParser(&post); err != nil {
		return helper.BadRequest(ctx, err)
	}

	if err := p.postService.Create(context.Background(), &post); err != nil {
		return helper.InternalServerError(ctx, err)
	}

	return helper.SuccessResponse(ctx, "post created successfully", post)
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}
