package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
)

func PostRoutes(app *fiber.App, postHandler *handlers.PostHandler) {
	app.Post("/posts", postHandler.CreatePostHandler)

}
