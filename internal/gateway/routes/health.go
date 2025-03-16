package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
)

func HealthCheckRoute(app *fiber.App, health *handlers.HealthHandler, auth helper.Auth) {
	app.Get("/v1/health", auth.BasicAuthentication, health.HealthCheckHandler)
}
