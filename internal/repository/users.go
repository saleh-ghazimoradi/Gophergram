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
	GetByID(ctx context.Context, id int64) (*service_modles.Users, error)
	CreateUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, id int64) error
	GetUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*service_modles.Users, error)
	UpdateUserInvitation(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error
	DeleteUserInvitation(ctx context.Context, tx *sql.Tx, id int64) error
	Delete(ctx context.Context, tx *sql.Tx, id int64) error
	GetByEmail(ctx context.Context, email string) (*service_modles.Users, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return u.db.BeginTx(ctx, nil)
}

func (u *userRepository) Create(ctx context.Context, tx *sql.Tx, user *service_modles.Users) error {
	query := `INSERT INTO users(username,password, email, role_id) VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4)) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()

	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(ctx, query, user.Username, user.Password.Hash, user.Email, role).Scan(&user.ID, &user.CreatedAt)
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

func (u *userRepository) GetByID(ctx context.Context, id int64) (*service_modles.Users, error) {
	query := `SELECT users.id, username, email, password, created_at, roles.*
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE users.id = $1 AND is_active = true`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	var user service_modles.Users
	err := u.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Role.ID, &user.Role.Name, &user.Role.Level, &user.Role.Description)
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

func (u *userRepository) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*service_modles.Users, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1 AND is_active = true`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.QueryTimeOut.Timeout)
	defer cancel()
	user := &service_modles.Users{}
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.Hash,
		&user.CreatedAt,
	)
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

func NewUserRepository(db *sql.DB) Users {
	return &userRepository{db: db}
}
