package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
)

type FeedHandler struct {
	postService service.PostService
}

func (f *FeedHandler) GetUserFeedHandler(ctx *fiber.Ctx) error {

	feed, err := f.postService.GetUserFeed(ctx.Context(), int64(4))
	if err != nil {
		return helper.InternalServerError(ctx, err)
	}

	return helper.SuccessResponse(ctx, fiber.StatusOK, "Success", feed)
}

func NewFeedHandler(postService service.PostService) *FeedHandler {
	return &FeedHandler{
		postService: postService,
	}
}
