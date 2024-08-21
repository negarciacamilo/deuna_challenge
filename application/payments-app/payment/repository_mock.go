package payment

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/domain/database"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) AddPayment(ctx *d.ContextInformation, payment *database.Payment) apierrors.ApiError {
	args := r.Called(ctx, payment)
	err := args.Get(0)
	if err != nil {
		return err.(apierrors.ApiError)
	}
	return nil
}

func (r *RepositoryMock) ChangePaymentStatus(ctx *d.ContextInformation, payment *database.Payment) apierrors.ApiError {
	args := r.Called(ctx, payment)
	err := args.Get(0)
	if err != nil {
		return err.(apierrors.ApiError)
	}
	return nil
}

func (r *RepositoryMock) GetAllPayments(ctx *d.ContextInformation) (*[]database.Payment, apierrors.ApiError) {
	args := r.Called(ctx)
	p := args.Get(0)
	err := args.Get(1)
	if err != nil {
		if p != nil {
			return p.(*[]database.Payment), err.(apierrors.ApiError)
		}
		return nil, err.(apierrors.ApiError)
	}
	return p.(*[]database.Payment), nil
}

func (r *RepositoryMock) GetCustomerPayments(ctx *d.ContextInformation, id uint64) (*[]database.Payment, apierrors.ApiError) {
	args := r.Called(ctx, id)
	p := args.Get(0)
	err := args.Get(1)
	if err != nil {
		if p != nil {
			return p.(*[]database.Payment), err.(apierrors.ApiError)
		}
		return nil, err.(apierrors.ApiError)
	}
	return p.(*[]database.Payment), nil
}

func (r *RepositoryMock) GetPaymentByID(ctx *d.ContextInformation, id uint64) (*database.Payment, apierrors.ApiError) {
	args := r.Called(ctx, id)
	p := args.Get(0)
	err := args.Get(1)
	if err != nil {
		if p != nil {
			return p.(*database.Payment), err.(apierrors.ApiError)
		}
		return nil, err.(apierrors.ApiError)
	}
	return p.(*database.Payment), nil
}
