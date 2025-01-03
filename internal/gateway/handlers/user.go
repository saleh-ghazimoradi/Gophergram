package handlers

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
)

type userHandler struct {
	userService service.UserService
}

func (u *userHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIdParam(r)
	if err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	user, err := u.userService.GetById(context.Background(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrsNotFound):
			helper.NotFoundResponse(w, r, err)
			return
		default:
			helper.InternalServerError(w, r, err)
			return
		}
	}

	if err = json.JSONResponse(w, http.StatusOK, user); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}
