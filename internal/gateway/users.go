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

// GetUserByID godoc
//
// @Summary Fetches a user profile
// @Description Fetches a user profile by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id	path	int true "id"
// @Success 200 {object} service_modles.Users
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /user/{id} [get]
func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}
}

// FollowUserHandler godoc
//
// @Summary Follow a user
// @Description Follows a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param userID	path	int true "id"
// @Success 204 {object} string "User followed"
// @Failure 400 {object} error "User payload missing"
// @Failure 404 {object} error "User not found"
// @Security ApiKeyAuth
// @Router /user/{id}/follow [put]
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

// UnfollowUserHandler gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"id"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/user/{id}/unfollow [put]
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
