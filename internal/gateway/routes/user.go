package routes

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
)

func User(router *httprouter.Router, db *sql.DB) {
	userRepository := repository.NewUserRepository(db, db)
	followRepository := repository.NewFollowerRepository(db, db)
	followService := service.NewFollowerService(followRepository)
	userService := service.NewUserService(userRepository, db)
	userHandler := handlers.NewUserHandler(userService, followService)
	middle := middlewares.NewMiddleware(nil, userService)

	userMiddleware := middle.UserContextMiddleware
	router.Handler(http.MethodGet, "/v1/users/:id", userMiddleware(http.HandlerFunc(userHandler.GetUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/follow", userMiddleware(http.HandlerFunc(userHandler.FollowUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/unfollow", userMiddleware(http.HandlerFunc(userHandler.UnFollowUserHandler)))
}
