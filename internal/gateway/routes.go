package gateway

import (
	"github.com/justinas/alice"
	"net/http"
)

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
	PostsContextMiddleware func(http.Handler) http.Handler
	UsersContextMiddleware func(http.Handler) http.Handler
}

func Routes(handler Handlers) http.Handler {
	mux := http.NewServeMux()

	standard := alice.New(recoverPanic, logRequest, commonHeaders)
	postChain := alice.New(handler.PostsContextMiddleware)
	userChain := alice.New(handler.UsersContextMiddleware)

	mux.HandleFunc("GET /v1/health", healthCheckHandler)
	mux.Handle("POST /v1/post", handler.CreatePostHandler)
	mux.Handle("GET /v1/post/{id}", postChain.Then(handler.GetPostHandler))
	mux.Handle("DELETE /v1/post/{id}", postChain.Then(handler.DeletePostHandler))
	mux.Handle("PATCH /v1/post/{id}", postChain.Then(handler.UpdatePostHandler))
	mux.Handle("GET /v1/user/{id}", userChain.Then(handler.GetUserHandler))
	mux.Handle("PUT /v1/user/{id}/follow", userChain.Then(handler.FollowUserHandler))
	mux.Handle("PUT /v1/user/{id}/unfollow", userChain.Then(handler.UnfollowUserHandler))
	mux.Handle("GET /v1/user/feed", handler.GetUserFeedHandler)

	return standard.Then(mux)
}
