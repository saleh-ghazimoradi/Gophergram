package repository

import (
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func exec(db *sql.DB, tx *sql.Tx) boil.ContextExecutor {
	if tx != nil {
		return tx
	}
	return db
}
