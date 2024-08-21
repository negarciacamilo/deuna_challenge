package http

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/negarciacamilo/deuna_challenge/application/context"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/bank"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/database"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/payment"
	"github.com/negarciacamilo/deuna_challenge/application/response"
	"net/http"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Logger())
	router.Use(logRequestHandler())
	router.Use(generateContext())
	router.NoRoute(noRouteHandler)
	router.Use(idempotencyKeyCheck())
	router.Use(AuthorizeClient())
	mapRoutes(router)
	return router
}

func mapRoutes(router *gin.Engine) {
	db := database.New()
	httpClient := resty.New()

	paymentsRepo := payment.NewRepository(db)
	bankRepo := bank.NewRepository(httpClient)

	paymentsService := payment.NewService(bankRepo, paymentsRepo)
	paymentsHandler := payment.NewHandler(paymentsService)

	router.POST("/pay", paymentsHandler.Pay)
	router.GET("/payments/:payment_id", paymentsHandler.GetPaymentByID)
	router.GET("/customers/:customer_id/payments", paymentsHandler.GetCustomerPayments)
	router.GET("/payments", paymentsHandler.GetAllPayments)
	router.PUT("/payments/:payment_id/refund", paymentsHandler.RefundPaymentByID)
	router.GET("/ping", ping)
}

func ping(c *gin.Context) {
	ctx := context.GetContextInformation(c)
	response.Respond(ctx, response.New(http.StatusOK, gin.H{"message": "pong"}), nil)
}
