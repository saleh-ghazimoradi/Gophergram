package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
)

func User(router *httprouter.Router, db *sql.DB) {
	userRepository := repository.NewUserRepository(db, db)
	userService := service.NewUserService(userRepository, db)
	userHandler := handlers.NewUserHandler(userService)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", userHandler.GetUserHandler)
}
