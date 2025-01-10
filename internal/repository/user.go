package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *service_models.User) error
	GetById(ctx context.Context, id int64) (*service_models.User, error)
	CreateUserInvitation(ctx context.Context, token string, exp time.Duration, id int64) error
	GetUserFromInvitation(ctx context.Context, token string) (*service_models.User, error)
	GetByEmail(ctx context.Context, email string) (*service_models.User, error)
	UpdateUserInvitation(ctx context.Context, user *service_models.User) error
	DeleteUserInvitation(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	WithTx(tx *sql.Tx) UserRepository
}

type userRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (u *userRepository) Create(ctx context.Context, user *service_models.User) error {

	query := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	args := []any{user.Username, user.Password.Hash, user.Email}
	if err := u.dbWrite.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (u *userRepository) GetById(ctx context.Context, id int64) (*service_models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1 AND is_active = true`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()
	var user service_models.User
	err := u.dbRead.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) CreateUserInvitation(ctx context.Context, token string, exp time.Duration, id int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	args := []any{token, id, time.Now().Add(exp)}
	if _, err := u.dbWrite.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetUserFromInvitation(ctx context.Context, token string) (*service_models.User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u JOIN user_invitations ui ON u.id = ui.user_id WHERE ui.token = $1 AND ui.expiry > $2`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	user := &service_models.User{}

	if err := u.dbRead.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) DeleteUserInvitation(ctx context.Context, id int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	_, err := u.dbWrite.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) UpdateUserInvitation(ctx context.Context, user *service_models.User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	_, err := u.dbWrite.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()
	_, err := u.dbWrite.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*service_models.User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	user := &service_models.User{}
	err := u.dbRead.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *userRepository) WithTx(tx *sql.Tx) UserRepository {
	return &userRepository{
		dbRead:  u.dbRead,
		dbWrite: u.dbWrite,
		tx:      tx,
	}
}

func NewUserRepository(dbRead, dbWrite *sql.DB) UserRepository {
	return &userRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
