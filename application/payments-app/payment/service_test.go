package payment

import (
	"github.com/google/uuid"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/domain/database"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/bank"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/defines"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPay(t *testing.T) {

}

func TestPayHappyPath(t *testing.T) {
	bankMock := new(bank.RepositoryMock)
	paymentRepoMock := new(RepositoryMock)

	id, _ := uuid.NewV7()
	i := id.String()
	bankMock.On("Pay", mock.Anything, mock.Anything).Return(i, nil)
	paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(nil)

	paymentService := NewService(bankMock, paymentRepoMock)

	resp, err := paymentService.Pay(d.TestContext(), domain.PaymentRequest{})

	bankMock.AssertCalled(t, "Pay", mock.Anything, mock.Anything)
	paymentRepoMock.AssertCalled(t, "AddPayment", mock.Anything, mock.Anything)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Response().(*database.Payment))
	require.Equal(t, i, *resp.Response().(*database.Payment).OperationID)
}

func TestPayBankError(t *testing.T) {
	bankMock := new(bank.RepositoryMock)
	paymentRepoMock := new(RepositoryMock)

	bankMock.On("Pay", mock.Anything, mock.Anything).Return("", apierrors.NewBadRequestApiError("invalid card"))
	paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(nil)

	paymentService := NewService(bankMock, paymentRepoMock)

	resp, err := paymentService.Pay(d.TestContext(), domain.PaymentRequest{})

	require.Equal(t, defines.REJECTED_STATUS, resp.Response().(*database.Payment).Status)
	require.Nil(t, err)
}

func TestPayRepositoryError(t *testing.T) {
	bankMock := new(bank.RepositoryMock)
	paymentRepoMock := new(RepositoryMock)

	id, _ := uuid.NewV7()
	i := id.String()
	bankMock.On("Pay", mock.Anything, mock.Anything).Return(i, nil)
	bankMock.On("ReverseOperation", mock.Anything, i).Return(nil)
	paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(apierrors.NewBadRequestApiError("invalid card"))
	paymentRepoMock.On("ChangePaymentStatus", mock.Anything, mock.Anything).Return(nil)

	paymentService := NewService(bankMock, paymentRepoMock)

	resp, err := paymentService.Pay(d.TestContext(), domain.PaymentRequest{})

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, defines.REVERSAL_STATUS, resp.Response().(*database.Payment).Status)
}
