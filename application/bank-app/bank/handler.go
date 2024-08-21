package bank

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/context"
	"github.com/negarciacamilo/deuna_challenge/application/defines"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	d "github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/negarciacamilo/deuna_challenge/application/response"
	"github.com/spf13/viper"
)

type Handler interface {
	Pay(c *gin.Context)
	PerformReversal(c *gin.Context)
}

type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) Pay(c *gin.Context) {
	ctx := context.GetContextInformation(c)
	var clientHasEnoughBalanceRequest d.PaymentRequest

	apierr := context.ShouldBindJSON(ctx, &clientHasEnoughBalanceRequest)
	if apierr != nil {
		response.Respond(ctx, nil, apierr)
		return
	}

	cardHashIsValid := viper.GetBool("CARD_HASH_IS_VALID")
	enoughBalance := viper.GetBool("CLIENT_HAS_ENOUGH_BALANCE")
	exceededLimit := viper.GetBool("CLIENT_HAS_EXCEEDED_LIMIT")
	bankTxFailed := viper.GetBool("BANK_TX_FAILED")

	if !cardHashIsValid {
		response.Respond(ctx, nil, apierrors.NewBadRequestApiError(defines.INVALID_CARD_HASH))
		return
	}

	if !enoughBalance {
		response.Respond(ctx, nil, apierrors.NewBadRequestApiError(defines.CLIENT_INVALID_BALANCE))
		return
	}

	if exceededLimit {
		response.Respond(ctx, nil, apierrors.NewBadRequestApiError(defines.CLIENT_HAS_EXCEEDED_LIMIT))
		return
	}

	if bankTxFailed {
		response.Respond(ctx, nil, apierrors.NewInternalServerApiError(defines.BANK_TX_FAILED, errors.New("the operation couldn't be stored")))
		return
	}

	id, _ := uuid.NewV7()
	response.Respond(ctx, response.New(200, domain.BankResponse{PaymentID: id.String()}), nil)
}

func (h *handler) PerformReversal(c *gin.Context) {
	ctx := context.GetContextInformation(c)
	response.Respond(ctx, response.New(200, nil), nil)
}
