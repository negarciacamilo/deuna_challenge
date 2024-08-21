package bank

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/defines"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (*string, apierrors.ApiError) {
	args := r.Called(ctx, payment)
	id := args.String(0)
	err := args.Get(1)
	if err != nil {
		if id != "" {
			return &id, err.(apierrors.ApiError)
		}
		return nil, err.(apierrors.ApiError)
	}
	return &id, nil
}

func (r *RepositoryMock) ReverseOperation(ctx *d.ContextInformation, operationID string) apierrors.ApiError {
	args := r.Called(ctx, operationID)
	err := args.Get(0)
	if err != nil {
		return err.(apierrors.ApiError)
	}
	return nil
}

func (r *RepositoryMock) ParseAPIError(apierr apierrors.ApiError) string {
	switch apierr.Message() {
	case defines.INVALID_CARD_HASH:
		return "1011"
	case defines.CLIENT_INVALID_BALANCE:
		return "1016"
	case defines.CLIENT_HAS_EXCEEDED_LIMIT:
		return "2011"
	default:
		return "9999"
	}
}

func (r *RepositoryMock) RefundPayment(ctx *d.ContextInformation, operationID string) apierrors.ApiError {
	args := r.Called(ctx, operationID)
	err := args.Get(0)
	if err != nil {
		return err.(apierrors.ApiError)
	}
	return nil
}
