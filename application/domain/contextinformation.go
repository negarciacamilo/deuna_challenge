package domain

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type ContextInformation struct {
	RequestInfo *RequestInfo
	GinContext  *gin.Context
}

type AuthenticatedUser struct {
	ClientID uint64 `json:"client_id"`
}

type RequestInfo struct {
	RequestID         string
	AuthenticatedUser *AuthenticatedUser
	IdempotencyKey    *string
}

func (c *ContextInformation) GetCtx() context.Context {
	return c.GinContext
}
