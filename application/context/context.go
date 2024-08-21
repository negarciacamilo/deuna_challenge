package context

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	"io"
	"strings"
)

const (
	RequestInfoKey        = "requestInfo"
	ContextInformationKey = "contextInformation"
)

func GetContextInformation(c *gin.Context) *domain.ContextInformation {
	if c != nil {
		if c.Keys != nil && c.Keys[ContextInformationKey] != nil {
			ctx, ok := c.Keys[ContextInformationKey].(*domain.ContextInformation)
			ctx.GinContext = c
			if ok {
				return ctx
			}
		}
	}
	return nil
}

func ShouldBindJSON(c *domain.ContextInformation, i interface{}) apierrors.ApiError {
	if err := c.GinContext.ShouldBindJSON(i); err != nil {
		var apierr apierrors.ApiError
		if errors.Is(err, io.EOF) {
			apierr = apierrors.NewBadRequestApiError("empty body, a body was expected")
		} else {
			apierr = apierrors.NewBadRequestApiError(err.Error())
		}
		logger.Error(apierr.Message(), strings.ToLower(strings.ReplaceAll(logger.GetCallerFunctionName(), ".", "-")), err, c)
		return apierr
	}
	return nil
}
