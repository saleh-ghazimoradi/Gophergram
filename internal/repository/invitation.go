package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/Gophergram/internal/repository/boiler_models"
	"github.com/saleh-ghazimoradi/Gophergram/sLogger"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type InvitationRepository interface {
	CreateUserInvitation(ctx context.Context, userInvitation *boiler_models.UserInvitation) error
	withTx(tx *sql.Tx) InvitationRepository
}

type invitationRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
	tx      *sql.Tx
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
