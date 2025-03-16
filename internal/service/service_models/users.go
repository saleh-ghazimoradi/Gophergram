package service_models

import (
	"time"
)

type Users struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
	Token     string    `json:"token"`
}

type UserWithToken struct {
	Users Users
	Token string `json:"token"`
}
