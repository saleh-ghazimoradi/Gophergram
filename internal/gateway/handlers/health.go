package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"net/http"
)

type HealthHandler struct{}

func (h *HealthHandler) HealthCheckHandler(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"env":     config.AppConfig.ServerConfig.Env,
		"version": config.AppConfig.ServerConfig.Version,
	})
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
