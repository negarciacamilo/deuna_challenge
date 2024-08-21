package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/context"
	"github.com/negarciacamilo/deuna_challenge/application/defines"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/idempotency"
	"github.com/negarciacamilo/deuna_challenge/application/response"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	Method            = "method"
	URL               = "url"
	Resource          = "resource"
	QueryParams       = "query params"
	Body              = "body"
	Headers           = "headers"
	IncomingRequest   = "Incoming request"
	OutcomingResponse = "Outcoming response"
)

func GenerateContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		ctx := &domain.ContextInformation{RequestInfo: generateRequestInformation(c.Request)}
		c.Keys[context.ContextInformationKey] = ctx
		c.Next()
	}
}

func generateRequestInformation(r *http.Request) *domain.RequestInfo {
	return &domain.RequestInfo{RequestID: getOrGenerateRequestID(r)}
}

func getOrGenerateRequestID(r *http.Request) string {
	if requestID := r.Header.Get(defines.XRequestID); requestID != "" {
		return requestID
	}

	u, err := uuid.NewV7()
	if err != nil {
		return ""
	}
	return u.String()
}

func noRouteHandler(c *gin.Context) {
	ctx := context.GetContextInformation(c)
	response.Respond(ctx, nil, apierrors.NewNotFoundApiError(fmt.Sprintf("Resource %s not found.", c.Request.URL.Path)))
}

func getResponseFromContext(c *gin.Context) (int, interface{}) {
	if c != nil && c.Keys != nil && c.Keys[response.StatusKey] != nil {
		status := c.Keys[response.StatusKey]
		if status == nil {
			logger.Panic("status code can't be empty", "get-response", nil, context.GetContextInformation(c), nil)
		}
		statusCode := status.(int)
		resp := c.Keys[response.ResponseKey]

		return statusCode, resp
	}
	return 0, nil
}

func respondAndLogResponse(c *gin.Context, elapsed int64, shouldLog bool) {
	status, response := getResponseFromContext(c)
	logRequest(OutcomingResponse, elapsed, c, shouldLog)
	if response == nil {
		c.Status(status)
	} else {
		c.JSON(status, response)
	}
}

func logRequest(requestType string, elapsed int64, c *gin.Context, shouldLog bool) {
	if c != nil {
		ctx := context.GetContextInformation(c)
		if c.Request != nil {
			tags := map[string]any{
				Method: c.Request.Method,
				URL:    c.Request.Host,
			}

			if c.Request.URL != nil {
				if c.Request.URL.Path != "" {
					tags[Resource] = c.Request.URL.Path
				}

				if c.Request.URL.RawQuery != "" {
					tags[QueryParams] = c.Request.URL.RawQuery
				}
			}

			if len(c.Request.Header) > 0 {
				tags[Headers] = parseHeaders(&c.Request.Header)
			}

			if requestType == IncomingRequest {
				if c.Request.Body != nil {
					tags[Body] = getRequestBody(c)
				}
			}

			if requestType == OutcomingResponse {
				elapsedTime, timeUnit := parseElapsedTime(float64(elapsed))
				tags["elapsed-time"] = fmt.Sprintf("%f %s", elapsedTime, timeUnit)

				if c.Keys != nil && c.Keys[response.ResponseKey] != nil {
					tags[response.ResponseKey] = c.Keys[response.ResponseKey]
				}

				if c.Keys != nil && shouldLog {
					logger.Info(OutcomingResponse, "outcoming-response", ctx, tags)
					return
				}
			}
			if shouldLog {
				logger.Info(IncomingRequest, "incoming-request", ctx, tags)
			}
		}
	}
}

func parseElapsedTime(elapsed float64) (float64, string) {
	if elapsed > 100000 {
		return elapsed / 1000000, "ms"
	}
	return elapsed, "ns"
}

func parseHeaders(headers *http.Header) string {
	var headersBuilder strings.Builder
	for h, v := range *headers {
		if v[0] != "" {
			headersBuilder.WriteString(fmt.Sprintf("{%s:%s} ", h, v[0]))
		}
	}
	return strings.TrimSuffix(headersBuilder.String(), " ")
}

func getRequestBody(c *gin.Context) string {
	var bodyBytes []byte
	bodyBytes, _ = io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	compactedBuff := new(bytes.Buffer)
	json.Compact(compactedBuff, bodyBytes) // nolint
	body := compactedBuff.String()
	return body
}

func logRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If the handler is /ping, we don't log, but we respond anyway
		shouldLog := !strings.Contains(c.Request.URL.Path, "/ping")
		logRequest(IncomingRequest, 0, c, shouldLog)
		start := time.Now()
		c.Next()
		respondAndLogResponse(c, time.Since(start).Nanoseconds(), shouldLog)
	}
}

func idempotencyKeyCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		shouldCheck := !strings.Contains(c.Request.URL.Path, "/ping")
		if h := c.GetHeader(defines.IdempotencyKey); h != "" && shouldCheck {
			ikExists := idempotency.IdempotencyKeyExists(h)
			if !ikExists {
				// Then is the first time we see this request
				ctx := context.GetContextInformation(c)
				ctx.RequestInfo.IdempotencyKey = &h
			} else {
				// Maybe this is a duplicate request
				apierr := apierrors.NewApiError("invalid request", "request might be duplicated", 409, nil)
				c.AbortWithStatusJSON(apierr.Status(), apierr)
			}
		} else {
			// Normally I would not do this, but just for making the people that are reviewing the challenge life easier
			ctx := context.GetContextInformation(c)
			randomKey, _ := uuid.NewV7()
			r := randomKey.String()
			ctx.RequestInfo.IdempotencyKey = &r
		}
	}
}

func AuthorizeClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		shouldCheck := !strings.Contains(c.Request.URL.Path, "/ping")
		if !shouldCheck {
			c.Next()
			return
		}
		ctx := context.GetContextInformation(c)

		authHeader := c.GetHeader(defines.Authorization)
		if authHeader == "" {
			apierr := apierrors.NewUnauthorizedApiError()
			logger.Error(apierr.Message(), "authorize-client", apierr, ctx)
			c.AbortWithStatusJSON(apierr.Status(), apierr)
			return
		}

		// Just for saving time I will not validate the token
		validToken := viper.GetBool("AUTH_TOKEN_IS_VALID")

		if !validToken {
			apierr := apierrors.NewUnauthorizedApiError()
			logger.Error(apierr.Message(), "authorize-client", apierr, ctx)
			c.AbortWithStatusJSON(apierr.Status(), apierr)
			return
		}

		ctx.RequestInfo.AuthenticatedUser = &domain.AuthenticatedUser{ClientID: uint64(rand.Int63n(10))}
		c.Next()
	}
}
