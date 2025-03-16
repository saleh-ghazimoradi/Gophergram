package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/config"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type InvitationRepository interface {
	CreateUserInvitation(ctx context.Context, userInvitation *boiler_models.UserInvitation) error
	UpdateUserInvitation(ctx context.Context, user *boiler_models.User) error
	DeleteUserInvitation(ctx context.Context, id int64) error
	withTx(tx *sql.Tx) InvitationRepository
}

type invitationRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
}

func (i *invitationRepository) UpdateUserInvitation(ctx context.Context, user *boiler_models.User) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	rowsAffected, err := user.Update(ctx, exec(i.dbWrite, i.tx), boil.Infer())
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (i *invitationRepository) DeleteUserInvitation(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, config.AppConfig.Database.Timeout)
	defer cancel()

	affected, err := boiler_models.UserInvitations(qm.Where("user_id = ?", userID)).DeleteAll(ctx, exec(i.dbWrite, i.tx))
	if err != nil {
		sLogger.SLogger.Error("Failed to delete user invitation", "user_id", userID, "err", err.Error())
		return err
	}
	if affected == 0 {
		sLogger.SLogger.Warn("No invitation found to delete", "user_id", userID)
		return sql.ErrNoRows
	}
	return nil
}

func (i *invitationRepository) CreateUserInvitation(ctx context.Context, userInvitation *boiler_models.UserInvitation) error {
	if err := userInvitation.Insert(ctx, exec(i.dbWrite, i.tx), boil.Infer()); err != nil {
		sLogger.SLogger.Error("failed to insert the user invitation", "err", err.Error())
		return err
	}
	return nil
}

func (i *invitationRepository) withTx(tx *sql.Tx) InvitationRepository {
	return &invitationRepository{
		dbRead:  i.dbRead,
		dbWrite: i.dbWrite,
		tx:      tx,
	}
}

func NewInvitationRepository(dbRead, dbWrite *sql.DB) InvitationRepository {
	return &invitationRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}
