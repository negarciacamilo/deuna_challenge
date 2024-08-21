package bank

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/defines"
	d "github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/spf13/viper"
)

type Repository interface {
	Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (*string, apierrors.ApiError)
	ReverseOperation(ctx *d.ContextInformation, paymentID string) apierrors.ApiError
	ParseAPIError(apierr apierrors.ApiError) string
}

type repository struct {
	httpClient *resty.Client
}

func NewRepository(httpClient *resty.Client) Repository {
	return &repository{
		httpClient: httpClient,
	}
}

func (r *repository) Pay(ctx *d.ContextInformation, payment domain.PaymentRequest) (*string, apierrors.ApiError) {
	baseUrl := viper.GetString("BANK_API_URL")
	url := fmt.Sprintf("%s/pay", baseUrl)

	res, err := r.httpClient.R().EnableTrace().SetBody(payment).Post(url)
	if err != nil {
		apierr := apierrors.NewInternalServerApiError("something happened paying", err)
		logger.Error(apierr.Message(), "bank-payment-request", apierr, ctx)
		return nil, apierr
	}

	if res.IsError() {
		var apierr apierrors.ApiError
		_ = json.Unmarshal(res.Body(), apierr)
		logger.Error(apierr.Message(), "bank-payment-request", apierr, ctx, map[string]any{"body": string(res.Body())})
		return nil, apierrors.NewApiError("can't perform the payment", apierr.Error(), res.StatusCode(), nil)
	}

	var bankResponse d.BankResponse
	_ = json.Unmarshal(res.Body(), &bankResponse)
	return &bankResponse.PaymentID, nil
}

func (r *repository) ReverseOperation(ctx *d.ContextInformation, paymentID string) apierrors.ApiError {
	baseUrl := viper.GetString("BANK_API_URL")
	url := fmt.Sprintf("%s/reversal/%d", baseUrl, paymentID)

	res, err := r.httpClient.R().EnableTrace().Put(url)
	if err != nil {
		apierr := apierrors.NewInternalServerApiError("something happened reversing", err)
		logger.Error(apierr.Message(), "reverse-operation", apierr, ctx, map[string]any{"paymentID": paymentID})
		return apierr
	}

	if res.IsError() {
		var apierr apierrors.ApiError
		_ = json.Unmarshal(res.Body(), apierr)
		logger.Error(apierr.Message(), "reverse-operation", apierr, ctx, map[string]any{"body": string(res.Body()), "paymentID": paymentID})
		return apierrors.NewApiError("can't perform the reversal", apierr.Error(), res.StatusCode(), nil)
	}

	return nil
}

func (r *repository) ParseAPIError(apierr apierrors.ApiError) string {
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
