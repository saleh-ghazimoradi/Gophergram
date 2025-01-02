package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/utils"
)

type UserService interface {
	Create(ctx context.Context, user *service_models.User) error
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

func NewUserService(userRepo repository.UserRepository, db *sql.DB) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}
