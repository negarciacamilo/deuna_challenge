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
	GetAllPayments(ctx *d.ContextInformation) (*[]dbd.Payment, apierrors.ApiError)
	GetPaymentByID(ctx *d.ContextInformation, id uint64) (*dbd.Payment, apierrors.ApiError)
	GetCustomerPayments(ctx *d.ContextInformation, id uint64) (*[]dbd.Payment, apierrors.ApiError)
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

func (r *repository) GetAllPayments(ctx *d.ContextInformation) (*[]dbd.Payment, apierrors.ApiError) {
	var payments []dbd.Payment
	err := r.db.GetDB().NewSelect().Model(&payments).
		Relation("Customer").
		Relation("Merchant").
		Relation("Bank").Scan(ctx.GetCtx())
	if err != nil {
		return nil, r.db.HandleDBError(ctx, "payments", database.Fetching, err)
	}

	if len(payments) == 0 {
		return nil, apierrors.NewNotFoundApiError("no payments found")
	}

	return &payments, nil
}

func (r *repository) GetPaymentByID(ctx *d.ContextInformation, id uint64) (*dbd.Payment, apierrors.ApiError) {
	var payment dbd.Payment
	err := r.db.GetDB().NewSelect().Model(&payment).Where("?TableAlias.id = ?", id).
		Relation("Customer").
		Relation("Merchant").
		Relation("Bank").Scan(ctx.GetCtx())
	if err != nil {
		return nil, r.db.HandleDBError(ctx, "payment", database.Fetching, err)
	}

	return &payment, nil
}

func (r *repository) GetCustomerPayments(ctx *d.ContextInformation, id uint64) (*[]dbd.Payment, apierrors.ApiError) {
	var payments []dbd.Payment
	err := r.db.GetDB().NewSelect().Model(&payments).
		Where("customer_id = ?", id).
		Relation("Customer").
		Relation("Merchant").
		Relation("Bank").Scan(ctx.GetCtx())
	if err != nil {
		return nil, r.db.HandleDBError(ctx, "payments", database.Fetching, err)
	}

	if len(payments) == 0 {
		return nil, apierrors.NewNotFoundApiError("no payments found")
	}

	return &payments, nil
}
