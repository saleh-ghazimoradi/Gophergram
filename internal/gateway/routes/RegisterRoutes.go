package routes

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

func RegisterRoutes(router *httprouter.Router, db *sql.DB, client *redis.Client) {
	health := handlers.NewHealthHandler()

	userRepo := repository.NewUserRepository(db, db)
	followRepo := repository.NewFollowerRepository(db, db)
	postRepo := repository.NewPostRepository(db, db)
	commentRepo := repository.NewCommentRepository(db, db)
	roleRepo := repository.NewRoleRepository(db, db)
	cacheRepository := repository.NewCacheRepository(client)

	userService := service.NewUserService(userRepo, db)
	followService := service.NewFollowerService(followRepo)
	postService := service.NewPostService(postRepo, db)
	commentService := service.NewCommentService(commentRepo)
	mailService := service.NewMailer(config.AppConfig.Mail.ApiKey, config.AppConfig.Mail.FromEmail)
	JWTAuthenticator := service.NewJWTAuthenticator(config.AppConfig.Authentication.Secret, config.AppConfig.Authentication.Aud, config.AppConfig.Authentication.Iss)
	roleService := service.NewRoleService(roleRepo)
	cacheService := service.NewCacheService(cacheRepository)

	middleware := middlewares.NewMiddleware(postService, userService, JWTAuthenticator, roleService, cacheService)

	feedHandler := handlers.NewFeedHandler(postService)
	userHandler := handlers.NewUserHandler(userService, followService)
	postHandler := handlers.NewPostHandler(postService, commentService)
	authHandler := handlers.NewAuthHandler(userService, mailService, JWTAuthenticator)

	registerHealthRoutes(router, health, middleware)
	registerUserRoutes(router, userHandler, middleware, feedHandler)
	registerPostRoutes(router, postHandler, middleware)
	registerAuthenticationRoutes(router, authHandler)

	docsURL := fmt.Sprintf("%s/swagger/doc.json", config.AppConfig.ServerConfig.Port)
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.Handler(httpSwagger.URL(docsURL)))
}
