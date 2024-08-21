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
	tests := []struct {
		name             string
		bankPayReturn    string
		bankPayError     error
		paymentAddReturn error
		expectedStatus   int
		expectedErr      error
		setupMocks       func(bankMock *bank.RepositoryMock, paymentRepoMock *RepositoryMock)
	}{
		{
			name:             "Happy path",
			bankPayReturn:    "some-unique-id",
			bankPayError:     nil,
			paymentAddReturn: nil,
			expectedStatus:   defines.APPROVED_STATUS,
			expectedErr:      nil,
			setupMocks: func(bankMock *bank.RepositoryMock, paymentRepoMock *RepositoryMock) {
				bankMock.On("Pay", mock.Anything, mock.Anything).Return("some-unique-id", nil)
				paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:             "Bank error",
			bankPayReturn:    "",
			bankPayError:     apierrors.NewBadRequestApiError("invalid card"),
			paymentAddReturn: nil,
			expectedStatus:   defines.REJECTED_STATUS,
			expectedErr:      nil,
			setupMocks: func(bankMock *bank.RepositoryMock, paymentRepoMock *RepositoryMock) {
				bankMock.On("Pay", mock.Anything, mock.Anything).Return("", apierrors.NewBadRequestApiError("invalid card"))
				paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:             "Repository error",
			bankPayReturn:    "some-unique-id",
			bankPayError:     nil,
			paymentAddReturn: apierrors.NewBadRequestApiError("invalid card"),
			expectedStatus:   defines.REVERSAL_STATUS,
			expectedErr:      nil,
			setupMocks: func(bankMock *bank.RepositoryMock, paymentRepoMock *RepositoryMock) {
				bankMock.On("Pay", mock.Anything, mock.Anything).Return("some-unique-id", nil)
				bankMock.On("ReverseOperation", mock.Anything, "some-unique-id").Return(nil)
				paymentRepoMock.On("AddPayment", mock.Anything, mock.Anything).Return(apierrors.NewBadRequestApiError("invalid card"))
				paymentRepoMock.On("ChangePaymentStatus", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bankMock := new(bank.RepositoryMock)
			paymentRepoMock := new(RepositoryMock)

			tt.setupMocks(bankMock, paymentRepoMock)

			paymentService := NewService(bankMock, paymentRepoMock)
			resp, err := paymentService.Pay(d.TestContext(), domain.PaymentRequest{})

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			if tt.expectedStatus != 0 {
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedStatus, resp.Response().(*database.Payment).Status)
			} else {
				require.NotNil(t, resp)
				require.NotNil(t, resp.Response().(*database.Payment))
				require.Equal(t, tt.bankPayReturn, *resp.Response().(*database.Payment).OperationID)
			}
		})
	}
}

func TestGetPaymentByID(t *testing.T) {

	tests := []struct {
		name            string
		expectedPayment database.Payment
		expectedErr     apierrors.ApiError
		setupMocks      func(paymentRepoMock *RepositoryMock)
	}{
		{
			name:            "Happy path",
			expectedPayment: database.Payment{},
			expectedErr:     nil,
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetPaymentByID", mock.Anything, mock.Anything).Return(&database.Payment{}, nil)
			},
		},
		{
			name:            "Repository error",
			expectedPayment: database.Payment{},
			expectedErr:     apierrors.NewBadRequestApiError("test"),
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetPaymentByID", mock.Anything, mock.Anything).Return(nil, apierrors.NewBadRequestApiError("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepoMock := new(RepositoryMock)

			tt.setupMocks(paymentRepoMock)

			paymentService := NewService(nil, paymentRepoMock)
			payment, err := paymentService.GetPaymentByID(d.TestContext(), 1)
			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.Equal(t, tt.expectedPayment, *payment.Response().(*database.Payment))
			}
		})
	}
}

func TestGetCustomerPayments(t *testing.T) {

	tests := []struct {
		name            string
		expectedPayment []database.Payment
		expectedErr     apierrors.ApiError
		setupMocks      func(paymentRepoMock *RepositoryMock)
	}{
		{
			name:            "Happy path",
			expectedPayment: []database.Payment{},
			expectedErr:     nil,
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetCustomerPayments", mock.Anything, mock.Anything).Return(&[]database.Payment{}, nil)
			},
		},
		{
			name:            "Repository error",
			expectedPayment: []database.Payment{},
			expectedErr:     apierrors.NewBadRequestApiError("test"),
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetCustomerPayments", mock.Anything, mock.Anything).Return(nil, apierrors.NewBadRequestApiError("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepoMock := new(RepositoryMock)

			tt.setupMocks(paymentRepoMock)

			paymentService := NewService(nil, paymentRepoMock)
			payments, err := paymentService.GetCustomerPayments(d.TestContext(), 1)
			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.Equal(t, tt.expectedPayment, *payments.Response().(*[]database.Payment))
			}
		})
	}
}

func TestGetAllPayments(t *testing.T) {

	tests := []struct {
		name            string
		expectedPayment []database.Payment
		expectedErr     apierrors.ApiError
		setupMocks      func(paymentRepoMock *RepositoryMock)
	}{
		{
			name:            "Happy path",
			expectedPayment: []database.Payment{},
			expectedErr:     nil,
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetAllPayments", mock.Anything).Return(&[]database.Payment{}, nil)
			},
		},
		{
			name:            "Repository error",
			expectedPayment: []database.Payment{},
			expectedErr:     apierrors.NewBadRequestApiError("test"),
			setupMocks: func(paymentRepoMock *RepositoryMock) {
				paymentRepoMock.On("GetAllPayments", mock.Anything).Return(nil, apierrors.NewBadRequestApiError("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepoMock := new(RepositoryMock)

			tt.setupMocks(paymentRepoMock)

			paymentService := NewService(nil, paymentRepoMock)
			payments, err := paymentService.GetAllPayments(d.TestContext())
			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.Equal(t, tt.expectedPayment, *payments.Response().(*[]database.Payment))
			}
		})
	}
}

func TestRefundPayment(t *testing.T) {
	id, _ := uuid.NewV7()
	i := id.String()
	tests := []struct {
		name        string
		expectedErr apierrors.ApiError
		setupMocks  func(paymentRepoMock *RepositoryMock, bankRepo *bank.RepositoryMock)
	}{
		{
			name:        "Happy path",
			expectedErr: nil,
			setupMocks: func(paymentRepoMock *RepositoryMock, bankRepo *bank.RepositoryMock) {
				paymentRepoMock.On("GetPaymentByID", mock.Anything, mock.Anything).Return(&database.Payment{Status: defines.APPROVED_STATUS, OperationID: &i}, nil)
				bankRepo.On("RefundPayment", mock.Anything, mock.Anything).Return(nil)
				paymentRepoMock.On("ChangePaymentStatus", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:        "Repository error",
			expectedErr: apierrors.NewBadRequestApiError("test"),
			setupMocks: func(paymentRepoMock *RepositoryMock, bankRepo *bank.RepositoryMock) {
				paymentRepoMock.On("GetPaymentByID", mock.Anything, mock.Anything).Return(nil, apierrors.NewBadRequestApiError("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentRepoMock := new(RepositoryMock)
			bankRepo := new(bank.RepositoryMock)
			tt.setupMocks(paymentRepoMock, bankRepo)

			paymentService := NewService(bankRepo, paymentRepoMock)
			payments, err := paymentService.RefundPayment(d.TestContext(), 1)
			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.Equal(t, defines.REFUNDED_STATUS, payments.Response().(*database.Payment).Status)
			}
		})
	}

}
