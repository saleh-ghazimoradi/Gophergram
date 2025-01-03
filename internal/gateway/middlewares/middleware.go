package middlewares

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
)

type customMiddleware struct {
	postService service.PostService
	userService service.UserService
}

func (m *customMiddleware) PostsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.ReadIdParam(r)
		if err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		post, err := m.postService.GetById(context.Background(), id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrsNotFound):
				helper.NotFoundResponse(w, r, err)
			default:
				helper.InternalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(context.Background(), handlers.PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *customMiddleware) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.ReadIdParam(r)
		if err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		user, err := m.userService.GetById(context.Background(), id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrsNotFound):
				helper.NotFoundResponse(w, r, err)
				return
			default:
				helper.InternalServerError(w, r, err)
				return
			}
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, handlers.UserCTX, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewMiddleware(postService service.PostService, userService service.UserService) *customMiddleware {
	return &customMiddleware{
		postService: postService,
		userService: userService,
	}
}
