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
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
	"strconv"
	"strings"
)

type CustomMiddleware struct {
	postService      service.PostService
	userService      service.UserService
	authService      service.Authenticator
	roleService      service.RoleService
	cacheService     service.CacheService
	rateLimitService service.RateLimitService
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

func (m *CustomMiddleware) CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

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

		user, err := m.getUser(context.Background(), userId)
		if err != nil {
			helper.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.UserCTX, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *CustomMiddleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *CustomMiddleware) getUser(ctx context.Context, id int64) (*service_models.User, error) {
	logger.Logger.Info("cache hit", "key", "user", "id", id)
	user, err := m.cacheService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		logger.Logger.Info("fetching from DB", "id", id)
		user, err = m.userService.GetById(ctx, id)
		if err != nil {
			return nil, err
		}
		if err = m.cacheService.Set(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}
func (m *CustomMiddleware) CheckPostOwnership(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := handlers.GetUserFromContext(r)
		post := handlers.GetPostFromCTX(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := m.checkRolePrecedence(context.Background(), user, requiredRole)
		if err != nil {
			helper.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			helper.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *CustomMiddleware) checkRolePrecedence(ctx context.Context, user *service_models.User, roleName string) (bool, error) {
	role, err := m.roleService.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

func (m *CustomMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		allowed, retryAfter, err := m.rateLimitService.IsAllowed(r.Context(), clientIP, config.AppConfig.Rate.Limit, config.AppConfig.Rate.Window)

		if err != nil && errors.Is(err, repository.ErrRateLimitExceeded) {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			helper.RateLimitExceededResponse(w, r, fmt.Sprintf("%v", config.AppConfig.Rate.Window))
			return
		} else if err != nil {
			helper.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			helper.RateLimitExceededResponse(w, r, fmt.Sprintf("%v", config.AppConfig.Rate.Window))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewMiddleware(postService service.PostService, userService service.UserService, authService service.Authenticator, roleService service.RoleService, cacheService service.CacheService, rateLimitService service.RateLimitService) *CustomMiddleware {
	return &CustomMiddleware{
		postService:      postService,
		userService:      userService,
		authService:      authService,
		roleService:      roleService,
		cacheService:     cacheService,
		rateLimitService: rateLimitService,
	}
}
