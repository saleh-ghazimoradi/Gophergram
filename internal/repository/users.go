package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *boiler_models.User) error
	GetById(ctx context.Context, id int64) (*boiler_models.User, error)
	CreateAndInvite(ctx context.Context, user *boiler_models.User) error
	GetUserFromInvitation(ctx context.Context, token string) (*boiler_models.User, error)
	Delete(ctx context.Context, id int64) error
	WithTx(tx *sql.Tx) UserRepository
}

type userRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (u *userRepository) Create(ctx context.Context, user *boiler_models.User) error {
	if err := user.Insert(ctx, exec(u.dbWrite, u.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the user", "err", err.Error())
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (u *userRepository) GetById(ctx context.Context, id int64) (*boiler_models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	user, err := boiler_models.FindUser(ctx, exec(u.dbRead, u.tx), id)
	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the user", "id", id)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return user, err
}

func (u *userRepository) CreateAndInvite(ctx context.Context, user *boiler_models.User) error {
	if err := user.Insert(ctx, exec(u.dbWrite, u.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the user", "err", err.Error())
		return err
	}
	return nil
}

func (u *userRepository) GetUserFromInvitation(ctx context.Context, token string) (*boiler_models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	user, err := boiler_models.Users(
		qm.InnerJoin("user_invitations ui ON users.id = ui.user_id"),
		qm.Where("ui.token = ?", token),
		qm.And("ui.expiry > ?", time.Now()),
	).One(ctx, exec(u.dbRead, u.tx))

	if err != nil {
		sLogger.SLogger.Error("failed to retrieve the user", "err", err.Error())
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrsNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	_, err := boiler_models.Users(
		boiler_models.UserWhere.ID.EQ(id),
	).DeleteAll(ctx, exec(u.dbRead, u.tx))

	if err != nil {
		sLogger.SLogger.Error("failed to delete the user", "id", id)
		return err
	}

	return nil
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
