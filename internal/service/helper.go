package service

import (
	"context"
	"database/sql"
)

func withTransaction[T any](ctx context.Context, db *sql.DB, txFunc func(tx *sql.Tx) (T, error)) (T, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return *new(T), err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	result, err := txFunc(tx)
	if err != nil {
		return *new(T), err
	}

	return result, nil
}
