package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return u.db.BeginTx(ctx, nil)
}

func (u *userRepository) Create(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error {
	query := `INSERT INTO users(username,password, email) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *sql.DB) Users {
	return &userRepository{db: db}
}
