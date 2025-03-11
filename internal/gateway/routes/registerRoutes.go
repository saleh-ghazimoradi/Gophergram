package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
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

	withTX := transaction.NewTransaction(db)
	postService := service.NewPostService(postRepository, withTX)
	commentService := service.NewCommentsService(commentRepository)
	userService := service.NewUserService(userRepository)
	followService := service.NewFollowService(followerRepository)

	postHandler := handlers.NewPostHandler(postService, commentService)
	userHandler := handlers.NewUserHandler(userService, followService)
	feedHandler := handlers.NewFeedHandler(postService)
	HealthCheckRoute(app, health)
	PostRoutes(app, postHandler)
	userRoutes(app, userHandler, feedHandler)
}
