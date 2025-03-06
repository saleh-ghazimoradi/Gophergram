package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
)

func HealthCheckRoute(app *fiber.App, health *handlers.HealthHandler) {
	app.Get("/v1/health", health.HealthCheckHandler)
}
