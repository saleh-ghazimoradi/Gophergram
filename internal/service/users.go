package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
)

type Users interface {
	Create(ctx context.Context, users *service_modles.Users) error
}

type userService struct {
	userRepo repository.Users
}

func (u *userService) Create(ctx context.Context, users *service_modles.Users) error {
	return u.userRepo.Create(ctx, users)
}

func NewServiceUser(repo repository.Users) Users {
	return &userService{
		userRepo: repo,
	}
}
