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
	tx, err := u.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if err := u.userRepo.Create(ctx, tx, users); err != nil {
		return err
	}
	return nil
}

func NewServiceUser(repo repository.Users) Users {
	return &userService{
		userRepo: repo,
	}
}
