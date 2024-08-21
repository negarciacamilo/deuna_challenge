package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/negarciacamilo/deuna_challenge/application/context"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/domain"
	"github.com/negarciacamilo/deuna_challenge/application/response"
)

type Handler interface {
	Pay(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Pay(c *gin.Context) {
	ctx := context.GetContextInformation(c)

	var paymentRequest domain.PaymentRequest

	apierr := context.ShouldBindJSON(ctx, &paymentRequest)
	if apierr != nil {
		response.Respond(ctx, nil, apierr)
		return
	}

	apierr = paymentRequest.Validate(ctx)
	if apierr != nil {
		response.Respond(ctx, nil, apierr)
		return
	}

	p, apierr := h.service.Pay(ctx, paymentRequest)
	response.Respond(ctx, p, apierr)

}
