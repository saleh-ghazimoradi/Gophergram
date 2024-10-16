package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"time"
)

type Users interface {
	Create(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*service_modles.Users, error)
	CreateUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, id int64) error
	GetUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*service_modles.Users, error)
	UpdateUserInvitation(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error
	DeleteUserInvitation(ctx context.Context, tx *sql.Tx, id int64) error
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

	err := tx.QueryRowContext(ctx, query, user.Username, user.Password.Hash, user.Email).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmails
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsernames
		default:
			return err
		}
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

func (u *userRepository) CreateUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, id int64) error {

	query := `INSERT INTO user_invitation (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, id, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*service_modles.Users, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u JOIN user_invitation ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	user := &service_modles.Users{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return user, nil

}

func (u *userRepository) UpdateUserInvitation(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) DeleteUserInvitation(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM user_invitation WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *sql.DB) Users {
	return &userRepository{db: db}
}
