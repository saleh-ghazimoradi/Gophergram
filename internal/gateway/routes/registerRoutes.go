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

	withTX := transaction.NewTransaction(db)
	postService := service.NewPostService(postRepository, withTX)
	commentService := service.NewCommentsService(commentRepository)

	postHandler := handlers.NewPostHandler(postService, commentService)

	HealthCheckRoute(app, health)
	PostRoutes(app, postHandler)
}
