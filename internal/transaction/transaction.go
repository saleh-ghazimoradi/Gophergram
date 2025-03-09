package transaction

import (
	"context"
	"database/sql"
	"github.com/friendsofgo/errors"
)

type Transaction interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	WithTx(ctx context.Context, fn func(*sql.Tx) error) error
}

type transaction struct {
	db *sql.DB
}

func (t *transaction) Begin(ctx context.Context) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, nil)
}

func (t *transaction) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (t *transaction) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (t *transaction) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := t.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			t.Rollback(tx)
			panic(p)
		} else if err != nil {
			t.Rollback(tx)
		} else {
			err = t.Commit(tx)
			if err != nil {
				t.Rollback(tx)
			}
		}
	}()

	err = fn(tx)
	return err
}

func NewTransaction(db *sql.DB) Transaction {
	return &transaction{
		db: db,
	}
}
