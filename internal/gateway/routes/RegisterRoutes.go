package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
)

func RegisterRoutes(router *httprouter.Router, db *sql.DB) {
	health := handlers.NewHealthHandler()

	userRepo := repository.NewUserRepository(db, db)
	followRepo := repository.NewFollowerRepository(db, db)
	postRepo := repository.NewPostRepository(db, db)
	commentRepo := repository.NewCommentRepository(db, db)

	userService := service.NewUserService(userRepo, db)
	followService := service.NewFollowerService(followRepo)
	postService := service.NewPostService(postRepo, db)
	commentService := service.NewCommentService(commentRepo)

	middleware := middlewares.NewMiddleware(postService, userService)

	feedHandler := handlers.NewFeedHandler(postService)
	userHandler := handlers.NewUserHandler(userService, followService)
	postHandler := handlers.NewPostHandler(postService, commentService)

	registerHealthRoutes(router, health)
	registerUserRoutes(router, userHandler, middleware, feedHandler)
	registerPostRoutes(router, postHandler, middleware)
}
