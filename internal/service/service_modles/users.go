package service_modles

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Users struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	CreatedAt time.Time `json:"create_at"`
	IsActive  bool      `json:"is_active"`
	RoleID    int64     `json:"role_id"`
	Role      Roles     `json:"role"`
}

type Password struct {
	Text *string
	Hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Text = &text
	p.Hash = hash
	return nil
}
