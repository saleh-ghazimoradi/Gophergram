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

type userHandler struct {
	userService     service.UserService
	followerService service.FollowerService
}

func (u *userHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := json.JSONResponse(w, http.StatusOK, user); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func (u *userHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (u *userHandler) UnFollowUserHandler(w http.ResponseWriter, r *http.Request) {
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

func NewUserHandler(userService service.UserService, followService service.FollowerService) *userHandler {
	return &userHandler{
		userService:     userService,
		followerService: followService,
	}
}
