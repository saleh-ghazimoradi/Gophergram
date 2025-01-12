package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerPostRoutes(router *httprouter.Router, handler *handlers.PostHandler, middleware *middlewares.CustomMiddleware) {
	authTokenMiddleware := middleware.AuthTokenMiddleware
	postMiddleware := middleware.PostsContextMiddleware
	checkOwnership := middleware.CheckPostOwnership
	rateLimitMiddleware := middleware.RateLimitMiddleware

	router.Handler(http.MethodPost, "/v1/posts", rateLimitMiddleware(authTokenMiddleware(http.HandlerFunc(handler.CreatePostHandler))))
	router.Handler(http.MethodGet, "/v1/posts/:id", rateLimitMiddleware(authTokenMiddleware(postMiddleware(http.HandlerFunc(handler.GetPostByIdHandler)))))
	router.Handler(http.MethodPatch, "/v1/posts/:id", rateLimitMiddleware(authTokenMiddleware(postMiddleware(checkOwnership("moderator", http.HandlerFunc(handler.UpdatePostHandler))))))
	router.Handler(http.MethodDelete, "/v1/posts/:id", rateLimitMiddleware(authTokenMiddleware(postMiddleware(checkOwnership("admin", http.HandlerFunc(handler.DeletePostHandler))))))
}
