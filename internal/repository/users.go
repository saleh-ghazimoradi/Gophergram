package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*service_modles.Users, error)
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

func (u *userRepository) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*service_modles.Users, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	var user service_modles.Users
	err := tx.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func NewUserRepository(db *sql.DB) Users {
	return &userRepository{db: db}
}
