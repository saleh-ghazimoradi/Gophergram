package transaction

import (
	"context"
	"database/sql"
)

type Transaction interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
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

func NewTransaction(db *sql.DB) Transaction {
	return &transaction{
		db: db,
	}
}

func WithTransaction(ctx context.Context, tm Transaction, fn func(tx *sql.Tx) error) error {
	tx, err := tm.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tm.Rollback(tx)
			panic(p)
		} else if err != nil {
			tm.Rollback(tx)
		} else {
			err = tm.Commit(tx)
		}
	}()

	err = fn(tx)
	return err
}
