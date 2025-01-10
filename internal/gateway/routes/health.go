package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerHealthRoutes(router *httprouter.Router, health *handlers.HealthHandler, middleware *middlewares.CustomMiddleware) {
	authMiddleware := middleware.BasicAuthentication
	router.Handler(http.MethodGet, "/v1/health", authMiddleware(http.HandlerFunc(health.Health)))
}
