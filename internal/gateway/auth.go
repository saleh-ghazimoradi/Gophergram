package gateway

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"net/http"
)

type UserWithToken struct {
	*service_modles.Users
	Token string `json:"token"`
}

type Auth struct {
	userService service.Users
}

// RegisterUserHandler godoc
//
// @Summary Register a user
// @Description Register a user
// @Tags authentication
// @Accept json
// @Produce json
// @Param payload body service_modles.RegisterUserPayLoad true "User credentials"
// @Success 201 {object} UserWithToken "User registered"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /authentication/user [post]
func (a *Auth) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_modles.RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := &service_modles.Users{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := a.userService.CreateAndInvite(ctx, user, hashToken, config.AppConfig.General.Mail.Exp)
	if err != nil {
		switch err {
		case repository.ErrDuplicateEmails:
			badRequestResponse(w, r, err)
		case repository.ErrDuplicateUsernames:
			badRequestResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	userWithToken := &UserWithToken{
		Users: user,
		Token: plainToken,
	}

	if err := jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		internalServerError(w, r, err)
	}
}

func NewAuth(userService service.Users) *Auth {
	return &Auth{userService: userService}
}
