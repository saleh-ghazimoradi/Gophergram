package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"log"
	"time"
)

type Users interface {
	Create(ctx context.Context, users *service_modles.Users) error
	GetByID(ctx context.Context, id int64) (*service_modles.Users, error)
	CreateAndInvite(ctx context.Context, user *service_modles.Users, token string, exp time.Duration) error
	ActivateUser(ctx context.Context, token string) error
	Delete(ctx context.Context, id int64) error
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

func (u *userService) GetByID(ctx context.Context, id int64) (*service_modles.Users, error) {
	tx, err := u.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error during transaction rollback: %v", rollbackErr)
			}
		}
	}()
	user, err := u.userRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if userErr := tx.Commit(); userErr != nil {
		return nil, userErr
	}
	return user, nil
}

func (u *userService) CreateAndInvite(ctx context.Context, user *service_modles.Users, token string, exp time.Duration) (err error) {
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

	if err = u.userRepo.Create(ctx, tx, user); err != nil {
		return err
	}

	if err = u.userRepo.CreateUserInvitation(ctx, tx, token, exp, user.ID); err != nil {
		return err
	}

	return nil
}

func (u *userService) ActivateUser(ctx context.Context, token string) error {
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

	user, err := u.userRepo.GetUserFromInvitation(ctx, tx, token)
	if err != nil {
		return err
	}
	user.IsActive = true

	if err := u.userRepo.UpdateUserInvitation(ctx, tx, user); err != nil {
		return err
	}
	if err := u.userRepo.DeleteUserInvitation(ctx, tx, user.ID); err != nil {
		return err
	}
	return nil
}

func (u *userService) Delete(ctx context.Context, id int64) error {
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

	if err := u.userRepo.Delete(ctx, tx, id); err != nil {
		return err
	}

	if err := u.userRepo.DeleteUserInvitation(ctx, tx, id); err != nil {
		return err
	}

	return nil
}

func NewServiceUser(repo repository.Users) Users {
	return &userService{
		userRepo: repo,
	}
}
