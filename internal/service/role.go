package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type RoleService interface {
	GetByName(ctx context.Context, name string) (*service_models.Role, error)
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func (s *roleService) GetByName(ctx context.Context, name string) (*service_models.Role, error) {
	return s.roleRepository.GetByName(ctx, name)
}

func NewRoleService(roleRepository repository.RoleRepository) RoleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}
