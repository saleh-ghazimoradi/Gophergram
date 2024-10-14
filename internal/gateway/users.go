package gateway

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

type User struct {
	userService   service.Users
	followService service.Follow
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}
}

func (u *User) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := GetUserFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
	}

	ctx := r.Context()

	if err := u.followService.Follow(ctx, followedUser.ID, payload.UserID); err != nil {
		switch err {
		case repository.ErrConflict:
			conflictResponse(w, r, err)
			return
		default:
			internalServerError(w, r, err)
			return
		}
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (u *User) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := GetUserFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := u.followService.Unfollow(ctx, unfollowedUser.ID, payload.UserID); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func (u *User) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := u.userService.GetByID(ctx, id)
		if err != nil {
			switch err {
			case repository.ErrNotFound:
				notFoundResponse(w, r, err)
			default:
				internalServerError(w, r, err)
			}
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) *service_modles.Users {
	user, _ := r.Context().Value(userCtx).(*service_modles.Users)
	return user
}

func NewUserHandler(userService service.Users, followService service.Follow) *User {
	return &User{
		userService:   userService,
		followService: followService,
	}
}
