package gateway

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
)

func Server() error {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024,
	})

	// Middlewares
	app.Use(recover.New()) // Prevents crashes from panics
	app.Use(logger.New())  // Logs incoming requests

	// Register routes
	routes.RegisterRoutes(app)

	sLogger.SLogger.Info("Starting server", "port", config.AppConfig.ServerConfig.Port)

	if err := app.Listen(":3000"); err != nil {
		sLogger.SLogger.Error("Failed to start server", "error", err)
	}

	return nil
}
