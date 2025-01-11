package middlewares

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type CustomMiddleware struct {
	postService service.PostService
	userService service.UserService
	authService service.Authenticator
}

func (m *CustomMiddleware) PostsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.ReadIdParam(r)
		if err != nil {
			helper.BadRequestResponse(w, r, err)
			return
		}

		post, err := m.postService.GetById(context.Background(), id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrsNotFound):
				helper.NotFoundResponse(w, r, err)
			default:
				helper.InternalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(context.Background(), handlers.PostCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *CustomMiddleware) BasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			helper.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header format must be Basic"))
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			helper.UnauthorizedBasicErrorResponse(w, r, err)
			return
		}

		username := config.AppConfig.Authentication.Username
		pass := config.AppConfig.Authentication.Password

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != pass {
			helper.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
		}

		next.ServeHTTP(w, r)
	})
}

func (m *CustomMiddleware) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header format must be Bearer"))
			return
		}

		token := parts[1]

		jwtToken, err := m.authService.ValidateToken(token)
		if err != nil {
			helper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			helper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		user, err := m.userService.GetById(context.Background(), userId)
		if err != nil {
			helper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.UserCTX, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewMiddleware(postService service.PostService, userService service.UserService, authService service.Authenticator) *CustomMiddleware {
	return &CustomMiddleware{
		postService: postService,
		userService: userService,
		authService: authService,
	}
}
