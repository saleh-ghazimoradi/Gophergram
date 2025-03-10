package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
)

func userRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	v1 := app.Group("/v1")
	v1.Get("/users/:id", userHandler.UsersContextMiddleware, userHandler.GetUserHandler)
	v1.Put("/users/:id/follow", userHandler.UsersContextMiddleware, userHandler.FollowUserHandler)
	v1.Put("/users/:id/unfollow", userHandler.UsersContextMiddleware, userHandler.UnfollowUserHandler)
}
