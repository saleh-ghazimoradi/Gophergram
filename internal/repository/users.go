package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type UserRepository interface {
	Create(ctx context.Context, user *boiler_models.User) error
	WithTx(tx *sql.Tx) UserRepository
}

type userRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (u *userRepository) Create(ctx context.Context, user *boiler_models.User) error {
	if err := user.Insert(ctx, exec(u.dbWrite, u.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the user", err)
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
