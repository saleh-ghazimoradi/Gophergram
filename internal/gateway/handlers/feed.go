package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"strconv"
)

const (
	defaultLimit  = "10"
	defaultOffset = "0"
	defaultSearch = ""
)

type FeedHandler struct {
	postService service.PostService
}

func (f *FeedHandler) GetUserFeedHandler(ctx *fiber.Ctx) error {
	offset, err := strconv.Atoi(ctx.Query("offset", defaultOffset))
	if err != nil {
		return helper.BadRequest(ctx, err)
	}

	limit, err := strconv.Atoi(ctx.Query("limit", defaultLimit))
	if err != nil {
		return helper.BadRequest(ctx, err)
	}

	search := ctx.Query("search", defaultSearch)

	feed, err := f.postService.GetUserFeed(ctx.Context(), int64(4), offset, limit, search)
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
