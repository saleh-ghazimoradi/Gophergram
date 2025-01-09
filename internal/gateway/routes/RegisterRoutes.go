package routes

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
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
	mailService := service.NewMailer(config.AppConfig.Mail.ApiKey, config.AppConfig.Mail.FromEmail)

	middleware := middlewares.NewMiddleware(postService, userService)

	feedHandler := handlers.NewFeedHandler(postService)
	userHandler := handlers.NewUserHandler(userService, followService)
	postHandler := handlers.NewPostHandler(postService, commentService)
	authHandler := handlers.NewAuthHandler(userService, mailService)

	registerHealthRoutes(router, health)
	registerUserRoutes(router, userHandler, authHandler, middleware, feedHandler)
	registerPostRoutes(router, postHandler, middleware)

	docsURL := fmt.Sprintf("%s/swagger/doc.json", config.AppConfig.ServerConfig.Port)
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.Handler(httpSwagger.URL(docsURL)))
}
