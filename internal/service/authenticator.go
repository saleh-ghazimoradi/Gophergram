package service

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Claims, error)
}

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Claims, error) {
	return nil, nil
}

func NewJWTAuthenticator(secret, aud, iss string) Authenticator {
	return &JWTAuthenticator{secret: secret, aud: aud, iss: iss}
}
