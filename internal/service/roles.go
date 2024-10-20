package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Roles interface {
	GetByName(ctx context.Context, name string) (*service_modles.Roles, error)
}

type roleService struct {
	roleRepo repository.Roles
}

func (r *roleService) GetByName(ctx context.Context, name string) (*service_modles.Roles, error) {
	tx, err := r.roleRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	roles, err := r.roleRepo.GetByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func NewRoleService(roleRepo repository.Roles) Roles {
	return &roleService{
		roleRepo: roleRepo,
	}
}
