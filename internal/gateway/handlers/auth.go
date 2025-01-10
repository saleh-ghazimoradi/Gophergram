package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/json"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
	"net/http"
	"time"
)

type AuthHandler struct {
	userService service.UserService
	mailService service.Mailer
	authService service.Authenticator
}

// RegisterUserHandler Register a user
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		service_models.RegisterUserPayload	true	"User credentials"
//	@Success		201		{string}	service_models.UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/v1/authentication/user [post]
func (a *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_models.RegisterUserPayload
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err := helper.Validate.Struct(payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	user := &service_models.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := a.userService.CreateAndInvite(context.Background(), user, hashToken, config.AppConfig.Mail.Exp); err != nil {
		switch err {
		case repository.ErrDuplicateEmail:
			helper.BadRequestResponse(w, r, err)
		case repository.ErrDuplicateUsername:
			helper.BadRequestResponse(w, r, err)
		default:
			helper.InternalServerError(w, r, err)
		}
		return
	}

	userWithToken := &service_models.UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", config.AppConfig.Mail.FrontendURL, plainToken)
	isProdEnv := config.AppConfig.ServerConfig.Env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	fmt.Println(config.AppConfig.Mail.UserWelcomeTemplate)

	status, err := a.mailService.Send(config.AppConfig.Mail.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		logger.Logger.Error("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := a.userService.Delete(context.Background(), user.ID); err != nil {
			logger.Logger.Error("error deleting user", "error", err)
		}
		helper.InternalServerError(w, r, err)
		return
	}

	logger.Logger.Info("Email sent", "status code", status)
	if err := json.JSONResponse(w, http.StatusCreated, userWithToken); err != nil {
		helper.InternalServerError(w, r, err)
	}

}

// CreateTokenHandler Register a user
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		service_models.CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string									"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/v1/authentication/token [post]
func (a *AuthHandler) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload service_models.CreateUserTokenPayload
	if err := json.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	if err := helper.Validate.Struct(payload); err != nil {
		helper.BadRequestResponse(w, r, err)
		return
	}

	user, err := a.userService.GetByEmail(context.Background(), payload.Email)
	if err != nil {
		switch err {
		case repository.ErrsNotFound:
			helper.UnauthorizedErrorResponse(w, r, err)
		default:
			helper.InternalServerError(w, r, err)
		}
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(config.AppConfig.Mail.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.AppConfig.Authentication.Iss,
		"aud": config.AppConfig.Authentication.Aud,
	}

	token, err := a.authService.GenerateToken(claims)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusCreated, token); err != nil {
		helper.InternalServerError(w, r, err)
	}
}

func NewAuthHandler(userService service.UserService, mailService service.Mailer, authService service.Authenticator) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		mailService: mailService,
		authService: authService,
	}
}
