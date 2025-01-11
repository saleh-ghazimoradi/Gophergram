package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*service_models.Role, error)
	WithTx(tx *sql.Tx) RoleRepository
}

type roleRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*service_models.Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Context.ContextTimeout)
	defer cancel()

	role := &service_models.Role{}
	err := r.dbRead.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Level,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) WithTx(tx *sql.Tx) RoleRepository {
	return &roleRepository{
		dbRead:  r.dbRead,
		dbWrite: r.dbWrite,
		tx:      tx,
	}
}

func NewRoleRepository(dbRead, dbWrite *sql.DB) RoleRepository {
	return &roleRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
