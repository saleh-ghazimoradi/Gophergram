package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"net/http"
)

func registerAuthenticationRoutes(router *httprouter.Router, authHandler *handlers.AuthHandler) {
	router.HandlerFunc(http.MethodPost, "/v1/authentication/user", authHandler.RegisterUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/authentication/token", authHandler.CreateTokenHandler)
}
