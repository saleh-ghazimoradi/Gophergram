package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/transaction"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	health := handlers.NewHealthHandler()

	postRepository := repository.NewPostRepository(db, db)
	commentRepository := repository.NewCommentRepository(db, db)
	userRepository := repository.NewUserRepository(db, db)
	followerRepository := repository.NewFollowerRepository(db, db)
	invitationRepository := repository.NewInvitationRepository(db, db)
	withTX := transaction.NewTransaction(db)

	maileService := service.NewMailer(config.AppConfig.Mail.FromEmail, config.AppConfig.Mail.ApiKey)
	postService := service.NewPostService(postRepository, withTX)
	commentService := service.NewCommentsService(commentRepository)
	authService := helper.NewAuth(config.AppConfig.Authentication.Secret)
	userService := service.NewUserService(userRepository, invitationRepository, authService, withTX, maileService)
	followService := service.NewFollowService(followerRepository)

	postHandler := handlers.NewPostHandler(postService, commentService)
	userHandler := handlers.NewUserHandler(userService, followService)
	feedHandler := handlers.NewFeedHandler(postService)
	HealthCheckRoute(app, health)
	PostRoutes(app, postHandler)
	userRoutes(app, userHandler, feedHandler)
}
