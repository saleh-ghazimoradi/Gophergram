package gateway

import (
	"fmt"
	"github.com/justinas/alice"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

type CheckOwnershipFunc func(requiredRole string, next http.HandlerFunc) http.HandlerFunc

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
	CheckPostOwnership     CheckOwnershipFunc
}

func Routes(handler Handlers) http.Handler {
	mux := http.NewServeMux()

	standard := alice.New(recoverPanic, logRequest, commonHeaders)
	postChain := alice.New(handler.PostsContextMiddleware)
	authChain := alice.New(basicAuthentication())
	authTokenChain := alice.New(handler.AuthTokenMiddleware)

	docsURL := fmt.Sprintf("%s/swagger/doc.json", config.AppConfig.General.Listen)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", httpSwagger.Handler(httpSwagger.URL(docsURL))))

	mux.Handle("/v1/health", authChain.Then(http.HandlerFunc(healthCheckHandler)))
	mux.Handle("POST /v1/post", authTokenChain.Then(handler.CreatePostHandler))
	mux.Handle("GET /v1/post/{id}", authTokenChain.Then(postChain.Then(handler.GetPostHandler)))
	mux.Handle("DELETE /v1/post/{id}", authTokenChain.Then(postChain.Then(handler.CheckPostOwnership("admin", handler.DeletePostHandler))))
	mux.Handle("PATCH /v1/post/{id}", authTokenChain.Then(postChain.Then(handler.CheckPostOwnership("moderator", handler.UpdatePostHandler))))
	mux.Handle("GET /v1/user/{id}", authTokenChain.Then(handler.GetUserHandler))
	mux.Handle("PUT /v1/user/{id}/follow", authTokenChain.Then(handler.FollowUserHandler))
	mux.Handle("PUT /v1/user/{id}/unfollow", authTokenChain.Then(handler.UnfollowUserHandler))
	mux.Handle("GET /v1/user/feed", authTokenChain.Then(handler.GetUserFeedHandler))
	mux.Handle("POST /v1/authentication/user", handler.RegisterUserHandler)
	mux.Handle("POST /v1/authentication/token", handler.CreateTokenHandler)

	mux.Handle("PUT /v1/activate/{token}", handler.ActivateUserHandler)

	return standard.Then(mux)
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
