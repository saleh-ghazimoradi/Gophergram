package gateway

import (
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"net/http"
	"strconv"
)

type User struct {
	userService service.Users
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
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
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		internalServerError(w, r, err)
	}
}

func NewUserHandler(userService service.Users) *User {
	return &User{
		userService: userService,
	}
}
