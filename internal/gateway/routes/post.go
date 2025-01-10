package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerPostRoutes(router *httprouter.Router, handler *handlers.PostHandler, middleware *middlewares.CustomMiddleware) {
	postMiddleware := middleware.PostsContextMiddleware
	authTokenMiddleware := middleware.AuthTokenMiddleware
	router.Handler(http.MethodPost, "/v1/posts", authTokenMiddleware(http.HandlerFunc(handler.CreatePostHandler)))
	router.Handler(http.MethodGet, "/v1/posts/:id", authTokenMiddleware(postMiddleware(http.HandlerFunc(handler.GetPostByIdHandler))))
	router.Handler(http.MethodPatch, "/v1/posts/:id", authTokenMiddleware(postMiddleware(http.HandlerFunc(handler.UpdatePostHandler))))
	router.Handler(http.MethodDelete, "/v1/posts/:id", authTokenMiddleware(postMiddleware(http.HandlerFunc(handler.DeletePostHandler))))
}
