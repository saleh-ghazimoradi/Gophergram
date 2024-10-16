package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
)

type feedHandler struct {
	postService service.Posts
}

// GetUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]service_modles.PostWithMetaData
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/feed [get]
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
