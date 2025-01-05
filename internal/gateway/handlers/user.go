package handlers

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"net/http"
)

type UserKey string

const UserCTX UserKey = "user"

type UserHandler struct {
	userService     service.UserService
	followerService service.FollowerService
}

// GetUserHandler retrieves the current user from the context.
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	service_models.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/users/{id} [get]
func (u *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := json.JSONResponse(w, http.StatusOK, user); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

// FollowUserHandler allows a user to follow another user.
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/v1/users/{id}/follow [put]
func (u *UserHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromContext(r)

	var payload service_models.FollowUser
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err := u.followerService.Follow(context.Background(), followedUser.ID, payload.UserID); err != nil {
		switch {
		case errors.Is(err, repository.ErrsConflict):
			helper.ConflictResponse(w, r, err)
			return
		default:
			helper.InternalServerError(w, r, err)
			return
		}
	}

	if err := json.JSONResponse(w, http.StatusNoContent, nil); err != nil {
	}
}

// UnFollowUserHandler allows a user to unfollow another user.
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/v1/users/{id}/unfollow [put]
func (u *UserHandler) UnFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := getUserFromContext(r)

	var payload service_models.FollowUser
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err := u.followerService.Unfollow(context.Background(), unFollowedUser.ID, payload.UserID); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusNoContent, nil); err != nil {
	}
}

func getUserFromContext(r *http.Request) *service_models.User {
	user, _ := r.Context().Value(UserCTX).(*service_models.User)
	return user
}

func NewUserHandler(userService service.UserService, followService service.FollowerService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		followerService: followService,
	}
}
