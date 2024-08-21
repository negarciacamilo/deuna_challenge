package payment

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/negarciacamilo/deuna_challenge/application/response"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (response.Response, apierrors.ApiError) {
	args := s.Called(ctx, payment)
	resp := args.Get(0)
	err := args.Get(1)
	if resp != nil {
		return resp.(response.Response), nil
	}
	if err != nil {
		return nil, err.(apierrors.ApiError)
	}

	return resp.(response.Response), err.(apierrors.ApiError)
}

func (s *ServiceMock) GetPaymentByID(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError) {
	args := s.Called(ctx, id)
	resp := args.Get(0)
	err := args.Get(1)
	if resp != nil {
		return resp.(response.Response), nil
	}
	if err != nil {
		return nil, err.(apierrors.ApiError)
	}

	return resp.(response.Response), err.(apierrors.ApiError)
}

func (s *ServiceMock) GetCustomerPayments(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError) {
	args := s.Called(ctx, id)
	resp := args.Get(0)
	err := args.Get(1)
	if resp != nil {
		return resp.(response.Response), nil
	}
	if err != nil {
		return nil, err.(apierrors.ApiError)
	}

	return resp.(response.Response), err.(apierrors.ApiError)
}

func (s *ServiceMock) GetAllPayments(ctx *d.ContextInformation) (response.Response, apierrors.ApiError) {
	args := s.Called(ctx)
	resp := args.Get(0)
	err := args.Get(1)
	if resp != nil {
		return resp.(response.Response), nil
	}
	if err != nil {
		return nil, err.(apierrors.ApiError)
	}

	return resp.(response.Response), err.(apierrors.ApiError)
}

func (s *ServiceMock) RefundPayment(ctx *d.ContextInformation, paymentID uint64) (response.Response, apierrors.ApiError) {
	args := s.Called(ctx, paymentID)
	resp := args.Get(0)
	err := args.Get(1)
	if resp != nil {
		return resp.(response.Response), nil
	}
	if err != nil {
		return nil, err.(apierrors.ApiError)
	}

	return resp.(response.Response), err.(apierrors.ApiError)
}
