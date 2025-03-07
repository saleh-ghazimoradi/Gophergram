package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
)

func PostRoutes(app *fiber.App, postHandler *handlers.PostHandler) {
	v1 := app.Group("/v1")
	v1.Post("/posts", postHandler.CreatePostHandler)
	v1.Get("/posts/:id", postHandler.GetPostHandler)
	v1.Delete("/posts/:id", postHandler.DeletePostHandler)
	//v1.Patch("/posts/:id", postHandler.UpdatePostHandler)
}
