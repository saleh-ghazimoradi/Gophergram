package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/middlewares"
	"net/http"
)

func registerUserRoutes(router *httprouter.Router, user *handlers.UserHandler, middleware *middlewares.CustomMiddleware, feed *handlers.FeedHandler) {
	authTokenMiddleware := middleware.AuthTokenMiddleware
	rateLimitMiddleware := middleware.RateLimitMiddleware
	recoverPanic := middleware.RecoverPanic
	commonHeader := middleware.CommonHeaders
	router.Handler(http.MethodPut, "/v1/user/activate/:token", commonHeader(recoverPanic(rateLimitMiddleware(http.HandlerFunc(user.ActivateUserHandler)))))
	router.Handler(http.MethodGet, "/v1/users/:id", commonHeader(recoverPanic(rateLimitMiddleware(authTokenMiddleware(http.HandlerFunc(user.GetUserHandler))))))
	router.Handler(http.MethodPut, "/v1/users/:id/follow", commonHeader(recoverPanic(rateLimitMiddleware(authTokenMiddleware(http.HandlerFunc(user.FollowUserHandler))))))
	router.Handler(http.MethodPut, "/v1/users/:id/unfollow", commonHeader(recoverPanic(rateLimitMiddleware(authTokenMiddleware(http.HandlerFunc(user.UnFollowUserHandler))))))
	router.Handler(http.MethodGet, "/v1/user/feed", commonHeader(recoverPanic(rateLimitMiddleware(authTokenMiddleware(http.HandlerFunc(feed.GetUserFeedHandler))))))
}
