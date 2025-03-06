package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type UserService interface {
	Create(ctx context.Context, user *service_models.Users) error
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) Create(ctx context.Context, user *service_models.Users) error {
	return u.userRepository.Create(ctx, user)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
