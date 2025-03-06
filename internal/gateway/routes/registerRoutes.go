package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
)

func RegisterRoutes(app *fiber.App) {
	health := handlers.NewHealthHandler()

	HealthCheckRoute(app, health)
}
