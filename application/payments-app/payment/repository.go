package payment

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	dbd "github.com/negarciacamilo/deuna_challenge/application/domain/database"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/database"
)

type Repository interface {
	AddPayment(ctx *d.ContextInformation, payment *dbd.Payment) apierrors.ApiError
	ChangePaymentStatus(ctx *d.ContextInformation, payment dbd.Payment) apierrors.ApiError
}

type repository struct {
	db database.Database
}

func NewRepository(db database.Database) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddPayment(ctx *d.ContextInformation, payment *dbd.Payment) apierrors.ApiError {
	_, err := r.db.GetDB().NewInsert().Model(payment).Exec(ctx.GetCtx())
	if err != nil {
		return r.db.HandleDBError(ctx, "payments", database.Creating, err)
	}

	return nil
}

func (r *repository) ChangePaymentStatus(ctx *d.ContextInformation, payment dbd.Payment) apierrors.ApiError {
	_, err := r.db.GetDB().NewUpdate().Model(payment).Where("id = ?", payment.ID).Exec(ctx.GetCtx())
	if err != nil {
		return r.db.HandleDBError(ctx, "payments", database.Updating, err)
	}
	return nil
}
