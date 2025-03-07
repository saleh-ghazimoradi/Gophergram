package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	health := handlers.NewHealthHandler()

	postRepository := repository.NewPostRepository(db, db)
	commentRepository := repository.NewCommentRepository(db, db)

	postService := service.NewPostService(postRepository)
	commentService := service.NewCommentsService(commentRepository)

	postHandler := handlers.NewPostHandler(postService, commentService)

	HealthCheckRoute(app, health)
	PostRoutes(app, postHandler)
}
