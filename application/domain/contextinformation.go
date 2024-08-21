package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/negarciacamilo/deuna_challenge/application/environment"
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
	if !environment.IsDockerEnv() {
		return context.Background()
	}
	return c.GinContext
}

func TestContext() *ContextInformation {
	return &ContextInformation{
		RequestInfo: &RequestInfo{
			AuthenticatedUser: &AuthenticatedUser{
				ClientID: 1,
			},
		},
		GinContext: &gin.Context{},
	}
}
