package payment

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	dbd "github.com/negarciacamilo/deuna_challenge/application/domain/database"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/bank"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/defines"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/negarciacamilo/deuna_challenge/application/response"
	"net/http"
)

type Service interface {
	Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (response.Response, apierrors.ApiError)
	GetPaymentByID(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError)
	GetCustomerPayments(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError)
	GetAllPayments(ctx *d.ContextInformation) (response.Response, apierrors.ApiError)
}

type service struct {
	bankRepository    bank.Repository
	paymentRepository Repository
}

func NewService(bankRepository bank.Repository, paymentRepository Repository) Service {
	return &service{
		bankRepository:    bankRepository,
		paymentRepository: paymentRepository,
	}
}

func (s *service) Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (response.Response, apierrors.ApiError) {
	p := &dbd.Payment{
		Amount:     payment.Amount,
		CustomerID: ctx.RequestInfo.AuthenticatedUser.ClientID,
		MerchantID: payment.MerchantID,
		BankID:     payment.BankID,
		Status:     defines.APPROVED_STATUS,
		Code:       defines.APPROVE_CODE,
	}

	paymentID, apierr := s.bankRepository.Pay(ctx, payment)
	if apierr != nil {
		code := s.bankRepository.ParseAPIError(apierr)
		p.Code = code
		p.Status = defines.REJECTED_STATUS
	} else {
		p.PaymentID = paymentID
	}

	apierr = s.paymentRepository.AddPayment(ctx, p)
	if apierr != nil && p.Status == defines.APPROVED_STATUS {
		err := s.bankRepository.ReverseOperation(ctx, *paymentID)
		// Best effort to reverse the payment
		if err != nil {
			logger.Error("error reversing payment", "payment-service-pay", err, ctx)
		}
		p.Status = defines.REVERSAL_STATUS
		err = s.paymentRepository.ChangePaymentStatus(ctx, *p)
		if err != nil {
			logger.Error("error changing payment status", "payment-service-pay", err, ctx)
		}
		return nil, apierr
	}

	return response.New(http.StatusCreated, p), nil
}

func (s *service) GetPaymentByID(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError) {
	payments, apierr := s.paymentRepository.GetPaymentByID(ctx, id)
	if apierr != nil {
		return nil, apierr
	}

	return response.New(http.StatusOK, payments), nil
}

func (s *service) GetCustomerPayments(ctx *d.ContextInformation, id uint64) (response.Response, apierrors.ApiError) {
	payments, apierr := s.paymentRepository.GetCustomerPayments(ctx, id)
	if apierr != nil {
		return nil, apierr
	}

	return response.New(http.StatusOK, payments), nil
}

func (s *service) GetAllPayments(ctx *d.ContextInformation) (response.Response, apierrors.ApiError) {
	payments, apierr := s.paymentRepository.GetAllPayments(ctx)
	if apierr != nil {
		return nil, apierr
	}

	return response.New(http.StatusOK, payments), nil
}
