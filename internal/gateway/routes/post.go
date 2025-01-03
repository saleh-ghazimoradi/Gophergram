package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
)

func Posts(router *httprouter.Router, db *sql.DB) {
	postRepository := repository.NewPostRepository(db, db)
	commentRepository := repository.NewCommentRepository(db, db)
	postService := service.NewPostService(postRepository, db)
	commentService := service.NewCommentService(commentRepository)
	postHandler := handlers.NewPostHandler(postService, commentService)
	middle := middlewares.NewMiddleware(postService, nil)

	postMiddleware := middle.PostsContextMiddleware
	router.HandlerFunc(http.MethodPost, "/v1/posts", postHandler.CreatePostHandler)
	router.Handler(http.MethodGet, "/v1/posts/:id", postMiddleware(http.HandlerFunc(postHandler.GetPostByIdHandler)))
	router.Handler(http.MethodPatch, "/v1/posts/:id", postMiddleware(http.HandlerFunc(postHandler.UpdatePostHandler)))
	router.Handler(http.MethodDelete, "/v1/posts/:id", postMiddleware(http.HandlerFunc(postHandler.DeletePostHandler)))
}
