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
