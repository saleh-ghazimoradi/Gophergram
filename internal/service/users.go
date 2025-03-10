package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
)

type UserService interface {
	Create(ctx context.Context, input *dto.User) error
	GetById(ctx context.Context, id int64) (*service_models.Users, error)
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

func (u *userService) GetById(ctx context.Context, id int64) (*service_models.Users, error) {
	boilerUser, err := u.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &service_models.Users{
		ID:       boilerUser.ID,
		Username: boilerUser.Username,
		Email:    boilerUser.Email,
		//Password:  boilerUser.Password,
		CreatedAt: boilerUser.CreatedAt,
	}, nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
