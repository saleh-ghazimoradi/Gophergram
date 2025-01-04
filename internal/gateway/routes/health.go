package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"net/http"
)

func registerHealthRoutes(router *httprouter.Router, health *handlers.HealthHandler) {
	router.HandlerFunc(http.MethodGet, "/v1/health", health.Health)
}
