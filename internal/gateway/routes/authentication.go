package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerAuthenticationRoutes(router *httprouter.Router, authHandler *handlers.AuthHandler, middleware *middlewares.CustomMiddleware) {
	rateLimitMiddleware := middleware.RateLimitMiddleware
	router.Handler(http.MethodPost, "/v1/authentication/user", rateLimitMiddleware(http.HandlerFunc(authHandler.RegisterUserHandler)))
	router.Handler(http.MethodPost, "/v1/authentication/token", rateLimitMiddleware(http.HandlerFunc(authHandler.CreateTokenHandler)))
}
