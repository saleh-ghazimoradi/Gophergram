package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/transaction"
)

type UserService interface {
	Create(ctx context.Context, input *dto.RegisterUser) error
	GetById(ctx context.Context, id int64) (*service_models.Users, error)
}

type userService struct {
	userRepository       repository.UserRepository
	invitationRepository repository.InvitationRepository
	authentication       helper.Auth
	transaction          transaction.Transaction
}

func (u *userService) Create(ctx context.Context, input *dto.RegisterUser) error {
	hashedPassword, err := u.authentication.CreateHashedPassword(input.Password)
	if err != nil {
		return err
	}
	return u.transaction.WithTx(ctx, func(tx *sql.Tx) error {
		repoWithTx := u.userRepository.WithTx(tx)
		boilerUser := &boiler_models.User{
			Username: input.Username,
			Email:    input.Email,
			Password: hashedPassword,
		}

		if err = repoWithTx.CreateAndInvite(ctx, boilerUser); err != nil {
			return err
		}

		plainToken := uuid.New().String()
		hash := sha256.Sum256([]byte(plainToken))
		hashToken := hex.EncodeToString(hash[:])

		if err := u.invitationRepository.CreateUserInvitation(ctx, &boiler_models.UserInvitation{
			Token:  hashToken,
			UserID: boilerUser.ID,
		}); err != nil {
		}

		return nil
	})
}

func (u *userService) GetById(ctx context.Context, id int64) (*service_models.Users, error) {
	boilerUser, err := u.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &service_models.Users{
		ID:        boilerUser.ID,
		Username:  boilerUser.Username,
		Email:     boilerUser.Email,
		CreatedAt: boilerUser.CreatedAt,
	}, nil
}

func NewUserService(userRepository repository.UserRepository, invitationRepository repository.InvitationRepository, authentication helper.Auth, transaction transaction.Transaction) UserService {
	return &userService{
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
		authentication:       authentication,
		transaction:          transaction,
	}
}
