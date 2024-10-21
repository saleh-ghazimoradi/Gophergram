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
	return r.roleRepo.GetByName(ctx, name)
}

func NewRoleService(roleRepo repository.Roles) Roles {
	return &roleService{
		roleRepo: roleRepo,
	}
}
