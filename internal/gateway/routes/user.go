package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerUserRoutes(router *httprouter.Router, user *handlers.UserHandler, middleware *middlewares.CustomMiddleware, feed *handlers.FeedHandler) {
	authTokenMiddleware := middleware.AuthTokenMiddleware

	router.Handler(http.MethodGet, "/v1/users/:id", authTokenMiddleware(http.HandlerFunc(user.GetUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/follow", authTokenMiddleware(http.HandlerFunc(user.FollowUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/unfollow", authTokenMiddleware(http.HandlerFunc(user.UnFollowUserHandler)))

	router.Handler(http.MethodGet, "/v1/user/feed", authTokenMiddleware(http.HandlerFunc(feed.GetUserFeedHandler)))
	router.HandlerFunc(http.MethodPut, "/v1/user/activate/:token", user.ActivateUserHandler)
}
