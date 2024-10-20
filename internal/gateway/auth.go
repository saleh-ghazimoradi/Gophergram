package gateway

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
	"time"
)

type UserWithToken struct {
	*service_modles.Users
	Token string `json:"token"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=2,max=72"`
}

type Auth struct {
	userService   service.Users
	mailService   service.Mailer
	authenticator service.Authenticator
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
		Role: service_modles.Roles{
			Name: "user",
		},
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

	activationURL := fmt.Sprintf("%s/confirm/%s", config.AppConfig.General.FrontendURL, plainToken)
	isProdEnv := config.AppConfig.Env.Env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	status, err := a.mailService.Send(service.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		logger.Logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := a.userService.Delete(ctx, user.ID); err != nil {
			logger.Logger.Errorw("error deleting user", "error", err)
		}
		internalServerError(w, r, err)
		return
	}

	logger.Logger.Infow("Email sent", "status code", status)

	if err := jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		internalServerError(w, r, err)
	}
}

// CreateTokenHandler godoc
//
// @Summary Creates a token
// @Description Creates a token for a user
// @Tags authentication
// @Accept json
// @Produce json
// @Param payload body CreateUserTokenPayload true "User credentials"
// @Success 200 {object} string 	"Token"
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /authentication/token [post]
func (a *Auth) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user, err := a.userService.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			unauthorizedBasicErrorResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(config.AppConfig.General.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.AppConfig.General.Auth.Token.TokenHost,
		"aud": config.AppConfig.General.Auth.Token.TokenHost,
	}
	token, err := a.authenticator.GenerateToken(claims)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, token); err != nil {
		internalServerError(w, r, err)
	}
}

func NewAuth(userService service.Users, mailerService service.Mailer, authService service.Authenticator) *Auth {
	return &Auth{userService: userService, mailService: mailerService, authenticator: authService}
}
