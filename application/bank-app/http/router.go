package http

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/negarciacamilo/deuna_challenge/application/bank-app/bank"
	"github.com/negarciacamilo/deuna_challenge/application/context"
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
	mapRoutes(router)
	return router
}

func mapRoutes(router *gin.Engine) {
	handler := bank.NewHandler()

	router.GET("/ping", ping)
	router.POST("/pay", handler.Pay)
	router.PUT("/reversal/:paymentID", handler.PerformReversal)
}

func ping(c *gin.Context) {
	ctx := context.GetContextInformation(c)
	response.Respond(ctx, response.New(http.StatusOK, gin.H{"message": "pong"}), nil)
}
