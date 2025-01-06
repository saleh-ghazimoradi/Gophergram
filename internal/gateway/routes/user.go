package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerUserRoutes(router *httprouter.Router, user *handlers.UserHandler, auth *handlers.AuthHandler, middleware *middlewares.CustomMiddleware, feed *handlers.FeedHandler) {
	userMiddleware := middleware.UserContextMiddleware

	router.Handler(http.MethodGet, "/v1/users/:id", userMiddleware(http.HandlerFunc(user.GetUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/follow", userMiddleware(http.HandlerFunc(user.FollowUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/unfollow", userMiddleware(http.HandlerFunc(user.UnFollowUserHandler)))
	router.HandlerFunc(http.MethodGet, "/v1/user/feed", feed.GetUserFeedHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/authentication", auth.RegisterUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/user/activate/:token", user.ActivateUserHandler)
}
