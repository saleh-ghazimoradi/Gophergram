package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/dto"
	"github.com/saleh-ghazimoradi/Gophergram/internal/gateway/helper"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"github.com/saleh-ghazimoradi/Gophergram/internal/transaction"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"time"
)

type UserService interface {
	Create(ctx context.Context, input *dto.RegisterUser) error
	GetById(ctx context.Context, id int64) (*service_models.Users, error)
	GetUserFromInvitation(ctx context.Context, token string) (*service_models.Users, error)
	Activate(ctx context.Context, token string) error
}

type userService struct {
	userRepository       repository.UserRepository
	invitationRepository repository.InvitationRepository
	authentication       helper.Auth
	transaction          transaction.Transaction
	mailer               Mailer
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

		if err = u.invitationRepository.CreateUserInvitation(ctx, &boiler_models.UserInvitation{
			Token:  hashToken,
			UserID: boilerUser.ID,
			Expiry: time.Now().Add(24 * time.Hour),
		}); err != nil {
			return err
		}

		activationURL := fmt.Sprintf("%s/confirm/%s", config.AppConfig.Mail.FrontendURL, plainToken)
		isProdEnv := config.AppConfig.ServerConfig.Env == "production"
		vars := struct {
			Username      string
			ActivationURL string
		}{
			Username:      boilerUser.Username,
			ActivationURL: activationURL,
		}

		_, err = u.mailer.Send(config.AppConfig.Mail.UserWelcomeTemplate, boilerUser.Username, boilerUser.Email, vars, !isProdEnv)
		if err != nil {
			sLogger.SLogger.Error("error sending welcome email", "error", err)

			if err = u.userRepository.Delete(ctx, boilerUser.ID); err != nil {
				sLogger.SLogger.Error("error deleting user", "error", err)
				err = u.invitationRepository.DeleteUserInvitation(ctx, boilerUser.ID)
				if err != nil {
					sLogger.SLogger.Error("error deleting invitation", "error", err)
					return err
				}
			}
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
		IsActive:  boilerUser.IsActive,
	}, nil
}

func (u *userService) GetUserFromInvitation(ctx context.Context, token string) (*service_models.Users, error) {
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	boilerUser, err := u.userRepository.GetUserFromInvitation(ctx, hashToken)
	if err != nil {
		return nil, err
	}

	return &service_models.Users{
		ID:        boilerUser.ID,
		Username:  boilerUser.Username,
		Email:     boilerUser.Email,
		CreatedAt: boilerUser.CreatedAt,
		IsActive:  boilerUser.IsActive,
		Token:     hashToken,
	}, nil
}

func (u *userService) Activate(ctx context.Context, token string) error {
	return u.transaction.WithTx(ctx, func(tx *sql.Tx) error {
		user, err := u.userRepository.GetUserFromInvitation(ctx, token)
		if err != nil {
			return err
		}
		if user == nil {
			return sql.ErrNoRows
		}

		user.IsActive = true
		if err = u.invitationRepository.UpdateUserInvitation(ctx, user); err != nil {
			return err
		}

		if err = u.invitationRepository.DeleteUserInvitation(ctx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func NewUserService(userRepository repository.UserRepository, invitationRepository repository.InvitationRepository, authentication helper.Auth, transaction transaction.Transaction, mailer Mailer) UserService {
	return &userService{
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
		authentication:       authentication,
		transaction:          transaction,
		mailer:               mailer,
	}
}
