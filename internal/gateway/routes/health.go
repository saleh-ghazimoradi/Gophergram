package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"net/http"
)

func HealthCheck(router *httprouter.Router) {
	health := handlers.NewHealthHandler()
	router.HandlerFunc(http.MethodGet, "/v1/health", health.Health)
}
