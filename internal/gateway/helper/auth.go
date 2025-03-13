package helper

import (
	"github.com/friendsofgo/errors"
	"golang.org/x/crypto/bcrypt"
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

func NewAuth(secret string) Auth {
	return Auth{Secret: secret}
}
