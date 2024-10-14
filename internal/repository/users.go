package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Users interface {
	Create(ctx context.Context, user *service_modles.Users) error
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) Create(ctx context.Context, user *service_modles.Users) error {
	query := `INSERT INTO users(username,password, email) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	err := u.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *sql.DB) Users {
	return &userRepository{db: db}
}
