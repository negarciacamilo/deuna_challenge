package response

import (
	"github.com/gin-gonic/gin"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRespond(t *testing.T) {
	type respondTest struct {
		name       string
		response   Response
		apierr     apierrors.ApiError
		ctx        *domain.ContextInformation
		wantResp   interface{}
		wantStatus int
	}

	tests := []respondTest{
		{
			name:       "should return response",
			ctx:        &domain.ContextInformation{GinContext: &gin.Context{}},
			response:   New(http.StatusOK, "test"),
			apierr:     nil,
			wantResp:   "test",
			wantStatus: 200,
		},
		{
			name:       "should return apierr",
			ctx:        &domain.ContextInformation{GinContext: &gin.Context{}},
			response:   nil,
			apierr:     apierrors.NewBadRequestApiError("test"),
			wantStatus: 400,
			wantResp:   apierrors.NewBadRequestApiError("test"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Respond(tc.ctx, tc.response, tc.apierr)
			require.Equal(t, tc.wantResp, tc.ctx.GinContext.Keys[ResponseKey])
			require.Equal(t, tc.wantStatus, tc.ctx.GinContext.Keys[StatusKey])
		})
	}
}
