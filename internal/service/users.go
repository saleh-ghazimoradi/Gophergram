package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"time"
)

type Users interface {
	Create(ctx context.Context, users *service_modles.Users) error
	GetByID(ctx context.Context, id int64) (*service_modles.Users, error)
	CreateAndInvite(ctx context.Context, user *service_modles.Users, token string, exp time.Duration) error
	ActivateUser(ctx context.Context, token string) error
	Delete(ctx context.Context, id int64) error
	GetByEmail(ctx context.Context, email string) (*service_modles.Users, error)
}

type userService struct {
	userRepo  repository.Users
	cacheRepo repository.Cacher
	db        *sql.DB
}

func (u *userService) Create(ctx context.Context, users *service_modles.Users) error {
	_, err := withTransaction(ctx, u.db, func(tx *sql.Tx) (struct{}, error) {
		return struct{}{}, u.userRepo.Create(ctx, tx, users)
	})
	return err
}

func (u *userService) GetByID(ctx context.Context, id int64) (*service_modles.Users, error) {
	var user *service_modles.Users
	if !config.AppConfig.Database.Redis.Enabled {
		return u.userRepo.GetByID(ctx, id)
	}

	user, err := u.cacheRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = u.userRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	if err := u.cacheRepo.Set(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) CreateAndInvite(ctx context.Context, user *service_modles.Users, token string, exp time.Duration) error {
	_, err := withTransaction(ctx, u.db, func(tx *sql.Tx) (struct{}, error) {
		if err := u.userRepo.Create(ctx, tx, user); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, u.userRepo.CreateUserInvitation(ctx, tx, token, exp, user.ID)
	})
	return err
}

func (u *userService) ActivateUser(ctx context.Context, token string) error {
	_, err := withTransaction(ctx, u.db, func(tx *sql.Tx) (struct{}, error) {
		user, err := u.userRepo.GetUserFromInvitation(ctx, tx, token)
		if err != nil {
			return struct{}{}, err
		}
		user.IsActive = true

		if err := u.userRepo.UpdateUserInvitation(ctx, tx, user); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, u.userRepo.DeleteUserInvitation(ctx, tx, user.ID)
	})
	return err
}

func (u *userService) Delete(ctx context.Context, id int64) error {
	_, err := withTransaction(ctx, u.db, func(tx *sql.Tx) (struct{}, error) {
		if err := u.userRepo.Delete(ctx, tx, id); err != nil {
			return struct{}{}, err
		}
		return struct{}{}, u.userRepo.DeleteUserInvitation(ctx, tx, id)
	})
	return err
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*service_modles.Users, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func NewServiceUser(repo repository.Users, cacheRepo repository.Cacher, db *sql.DB) Users {
	return &userService{
		userRepo:  repo,
		cacheRepo: cacheRepo,
		db:        db,
	}
}
