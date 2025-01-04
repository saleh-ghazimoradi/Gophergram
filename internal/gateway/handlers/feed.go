package handlers

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"net/http"
)

type FeedHandler struct {
	postService service.PostService
}

func (f *FeedHandler) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	p := service_models.PaginatedFeedQuery{
		Limit:  config.AppConfig.Pagination.Limit,
		Offset: config.AppConfig.Pagination.Offset,
		Sort:   config.AppConfig.Pagination.Sort,
	}

	fmt.Printf("limit: %d, type: %T, offset: %d, type: %T\n", p.Limit, p.Limit, p.Offset, p.Offset)

	fq, err := p.Parse(r)
	if err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err = helper.Validate.Struct(fq); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	//user := getUserFromContext(r)
	//fmt.Println(user.ID)

	feed, err := f.postService.GetUserFeed(context.Background(), int64(10), fq)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err = json.JSONResponse(w, http.StatusOK, feed); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func NewFeedHandler(postService service.PostService) *FeedHandler {
	return &FeedHandler{
		postService: postService,
	}
}
