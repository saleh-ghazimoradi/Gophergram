package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Roles interface {
	GetByName(ctx context.Context, tx *sql.Tx, name string) (*service_modles.Roles, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type roleRepo struct {
	db *sql.DB
}

func (r *roleRepo) GetByName(ctx context.Context, tx *sql.Tx, name string) (*service_modles.Roles, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`

	role := &service_modles.Roles{}
	err := tx.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func NewRoleRepo(db *sql.DB) Roles {
	return &roleRepo{
		db: db,
	}
}
