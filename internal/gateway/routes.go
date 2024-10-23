package gateway

import (
	"context"
	"expvar"
	"fmt"
	"github.com/justinas/alice"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
	"time"
)

type CheckOwnershipFunc func(requiredRole string, next http.HandlerFunc) http.HandlerFunc

type RateLimit func(ctx context.Context, rateLimiter service.RateLimiter, limit int, window time.Duration) func(http.Handler) http.Handler

type Handlers struct {
	CreatePostHandler      http.HandlerFunc
	UpdatePostHandler      http.HandlerFunc
	DeletePostHandler      http.HandlerFunc
	GetPostHandler         http.HandlerFunc
	CreateUserHandler      http.HandlerFunc
	GetUserHandler         http.HandlerFunc
	GetUserFeedHandler     http.HandlerFunc
	UpdateUserHandler      http.HandlerFunc
	DeleteUserHandler      http.HandlerFunc
	FollowUserHandler      http.HandlerFunc
	UnfollowUserHandler    http.HandlerFunc
	RegisterUserHandler    http.HandlerFunc
	ActivateUserHandler    http.HandlerFunc
	CreateTokenHandler     http.HandlerFunc
	PostsContextMiddleware func(http.Handler) http.Handler
	AuthTokenMiddleware    func(http.Handler) http.Handler
	RateLimitMiddleware    RateLimit
	CheckPostOwnership     CheckOwnershipFunc
}

func Routes(handler Handlers, rateLimiter service.RateLimiter) http.Handler {
	mux := http.NewServeMux()
	ctx := context.Background()

	standard := alice.New(recoverPanic, logRequest, commonHeaders)
	postChain := alice.New(handler.PostsContextMiddleware)
	authChain := alice.New(basicAuthentication())
	authTokenChain := alice.New(handler.AuthTokenMiddleware)

	ratechain := alice.New(handler.RateLimitMiddleware(ctx, rateLimiter, config.AppConfig.Rate.Limit, config.AppConfig.Rate.Time))

	docsURL := fmt.Sprintf("%s/swagger/doc.json", config.AppConfig.General.Listen)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(httpSwagger.URL(docsURL))))

	mux.Handle("/v1/health", ratechain.Then(authChain.Then(http.HandlerFunc(healthCheckHandler))))
	mux.Handle("POST /v1/post", ratechain.Then(authTokenChain.Then(handler.CreatePostHandler)))
	mux.Handle("GET /v1/post/{id}", ratechain.Then(authTokenChain.Then(postChain.Then(handler.GetPostHandler))))
	mux.Handle("DELETE /v1/post/{id}", ratechain.Then(authTokenChain.Then(postChain.Then(handler.CheckPostOwnership("admin", handler.DeletePostHandler)))))
	mux.Handle("PATCH /v1/post/{id}", ratechain.Then(authTokenChain.Then(postChain.Then(handler.CheckPostOwnership("moderator", handler.UpdatePostHandler)))))
	mux.Handle("GET /v1/user/{id}", ratechain.Then(authTokenChain.Then(handler.GetUserHandler)))
	mux.Handle("PUT /v1/user/{id}/follow", ratechain.Then(authTokenChain.Then(handler.FollowUserHandler)))
	mux.Handle("PUT /v1/user/{id}/unfollow", ratechain.Then(authTokenChain.Then(handler.UnfollowUserHandler)))
	mux.Handle("GET /v1/user/feed", ratechain.Then(authTokenChain.Then(handler.GetUserFeedHandler)))
	mux.Handle("POST /v1/authentication/user", ratechain.Then(handler.RegisterUserHandler))
	mux.Handle("POST /v1/authentication/token", ratechain.Then(handler.CreateTokenHandler))

	mux.Handle("PUT /v1/activate/{token}", ratechain.Then(handler.ActivateUserHandler))

	mux.Handle("GET /debug/vars", expvar.Handler())

	return metrics(standard.Then(mux))

}

//func Routes(handler Handlers) http.Handler {
//	mux := http.NewServeMux()
//
//	standard := alice.New(recoverPanic, logRequest, commonHeaders)
//	postChain := alice.New(handler.PostsContextMiddleware)
//	userChain := alice.New(handler.UsersContextMiddleware)
//
//	docsURL := fmt.Sprintf("%s/swagger/doc.json", config.AppConfig.General.Listen)
//
//	mux.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(httpSwagger.URL(docsURL))))
//
//	mux.HandleFunc("GET /v1/health", healthCheckHandler)
//	mux.Handle("POST /v1/post", handler.CreatePostHandler)
//	mux.Handle("GET /v1/post/{id}", postChain.Then(handler.GetPostHandler))
//	mux.Handle("DELETE /v1/post/{id}", postChain.Then(handler.DeletePostHandler))
//	mux.Handle("PATCH /v1/post/{id}", postChain.Then(handler.UpdatePostHandler))
//	mux.Handle("GET /v1/user/{id}", userChain.Then(handler.GetUserHandler))
//	mux.Handle("PUT /v1/user/{id}/follow", userChain.Then(handler.FollowUserHandler))
//	mux.Handle("PUT /v1/user/{id}/unfollow", userChain.Then(handler.UnfollowUserHandler))
//	mux.Handle("GET /v1/user/feed", handler.GetUserFeedHandler)
//	mux.Handle("POST /v1/authentication/user", handler.RegisterUserHandler)
//	mux.Handle("PUT /v1/user/activate/{token}", handler.ActivateUserHandler)
//
//	return standard.Then(mux)
//}
