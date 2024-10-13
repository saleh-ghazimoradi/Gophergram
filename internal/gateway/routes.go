package gateway

import (
	"github.com/justinas/alice"
	"net/http"
)

type Handlers struct {
	CreatePostHandler http.HandlerFunc
	UpdatePostHandler http.HandlerFunc
	DeletePostHandler http.HandlerFunc
	GetPostHandler    http.HandlerFunc
	CreateUserHandler http.HandlerFunc
	GetUserHandler    http.HandlerFunc
	UpdateUserHandler http.HandlerFunc
	DeleteUserHandler http.HandlerFunc
}

func Routes(handler Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", healthCheckHandler)
	mux.HandleFunc("POST /v1/post", handler.CreatePostHandler)
	mux.HandleFunc("GET /v1/posts/{id}", handler.GetPostHandler)
	standard := alice.New(recoverPanic, logRequest, commonHeaders)

	return standard.Then(mux)
}
