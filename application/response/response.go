package response

import (
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
)

const (
	StatusKey   = "status"
	ResponseKey = "response"
)

type Response interface {
	Status() int
	Response() interface{}
}

type response struct {
	statusCode int
	response   interface{}
}

func (r *response) Status() int {
	return r.statusCode
}

func (r *response) Response() interface{} {
	return r.response
}

func New(status int, r interface{}) Response {
	return &response{statusCode: status, response: r}
}

func Respond(c *domain.ContextInformation, response Response, apierror apierrors.ApiError) {
	if c.GinContext.Keys == nil {
		c.GinContext.Keys = make(map[string]interface{})
	}
	if response != nil {
		c.GinContext.Keys[StatusKey] = response.Status()
		c.GinContext.Keys[ResponseKey] = response.Response()
	} else {
		c.GinContext.Keys[StatusKey] = apierror.Status()
		c.GinContext.Keys[ResponseKey] = apierror
	}
}
