package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
)

type UserService interface {
	Create(ctx context.Context, input *dto.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) Create(ctx context.Context, input *dto.User) error {
	userBoiler := &boiler_models.User{}

	if err := u.userRepository.Create(ctx, userBoiler); err != nil {
		return err
	}
	return nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
