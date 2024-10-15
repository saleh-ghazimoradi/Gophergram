package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
)

type feedHandler struct {
	postService service.Posts
}

func (f *feedHandler) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	p := service_modles.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := p.Parse(r)

	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := f.postService.GetUserFeed(ctx, int64(100), fq)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		internalServerError(w, r, err)
	}
}

func NewFeedHandler(postService service.Posts) *feedHandler {
	return &feedHandler{
		postService: postService,
	}
}
