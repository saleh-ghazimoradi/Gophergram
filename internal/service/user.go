package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user *service_models.User) error
	GetById(ctx context.Context, id int64) (*service_models.User, error)
	GetByEmail(ctx context.Context, email string) (*service_models.User, error)
	CreateAndInvite(ctx context.Context, user *service_models.User, token string, invitationExp time.Duration) error
	Delete(ctx context.Context, id int64) error
	Activate(ctx context.Context, token string) error
}

type userService struct {
	userRepo repository.UserRepository
	db       *sql.DB
}

func (u *userService) Create(ctx context.Context, user *service_models.User) error {
	return utils.WithTransaction(ctx, u.db, func(tx *sql.Tx) error {
		userRepoWithTx := u.userRepo.WithTx(tx)
		return userRepoWithTx.Create(ctx, user)
	})
}

func (u *userService) GetById(ctx context.Context, id int64) (*service_models.User, error) {
	return u.userRepo.GetById(ctx, id)
}

func (u *userService) CreateAndInvite(ctx context.Context, user *service_models.User, token string, invitationExp time.Duration) error {
	return utils.WithTransaction(ctx, u.db, func(tx *sql.Tx) error {
		userRepoWithTx := u.userRepo.WithTx(tx)
		if err := userRepoWithTx.Create(ctx, user); err != nil {
			return err
		}
		if err := userRepoWithTx.CreateUserInvitation(ctx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (u *userService) Activate(ctx context.Context, token string) error {
	return utils.WithTransaction(ctx, u.db, func(tx *sql.Tx) error {
		user, err := u.userRepo.GetUserFromInvitation(ctx, token)
		if err != nil {
			return err
		}
		user.IsActive = true
		if err = u.userRepo.UpdateUserInvitation(ctx, user); err != nil {
			return err
		}
		if err = u.userRepo.DeleteUserInvitation(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (u *userService) Delete(ctx context.Context, id int64) error {
	return utils.WithTransaction(ctx, u.db, func(tx *sql.Tx) error {
		if err := u.userRepo.Delete(ctx, id); err != nil {
			return err
		}
		if err := u.userRepo.DeleteUserInvitation(ctx, id); err != nil {
			return err
		}
		return nil
	})
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*service_models.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func NewUserService(userRepo repository.UserRepository, db *sql.DB) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}
