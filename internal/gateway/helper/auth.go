package helper

import (
	"encoding/base64"
	"github.com/friendsofgo/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Auth struct {
	Secret string
}

func (a *Auth) CreateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error hashing password: " + err.Error())
	}
	return string(hashedPassword), nil
}

func (a *Auth) BasicAuthentication(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return UnauthorizedBasicErrorResponse(ctx, errors.New("Authorization header required"))
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return UnauthorizedBasicErrorResponse(ctx, errors.New("authorization header format must be Basic"))
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return UnauthorizedBasicErrorResponse(ctx, err)
	}

	username := config.AppConfig.Authentication.Username
	pass := config.AppConfig.Authentication.Password

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 || creds[0] != username || creds[1] != pass {
		return UnauthorizedBasicErrorResponse(ctx, errors.New("invalid credentials"))
	}

	return ctx.Next()
}

func NewAuth(secret string) Auth {
	return Auth{Secret: secret}
}
