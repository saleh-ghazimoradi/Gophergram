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
	CreateAndInvite(ctx context.Context, user *service_models.User, token string, invitationExp time.Duration) error
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

func NewUserService(userRepo repository.UserRepository, db *sql.DB) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}
