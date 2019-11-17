package middleware

import (
	"github.com/google/uuid"
	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/valyala/fasthttp"
)

//easyjson:json
type jsonErrorWrap struct {
	Code    int               `json:"code,omitempty"`
	UUID    string            `json:"uuid,omitempty"`
	Message string            `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

//easyjson:json
type jsonError struct {
	Err jsonErrorWrap `json:"error"`
}

func Json(handler Endpoint) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		responseHeader := &ctx.Response.Header
		responseHeader.SetContentType("application/json; charset=utf-8")
		responseHeader.Set("Access-Control-Allow-Origin", "*")

		err := handler(ctx)
		statusCode := fasthttp.StatusOK

		if err != nil {
			jsonErr := &jsonError{
				Err: jsonErrorWrap{
					UUID: uuid.New().String(),
				},
			}
			plusErr, ok := err.(errutil.Error)
			if ok {
				jsonErr.Err.Code = plusErr.GetCode()
				jsonErr.Err.Message = err.Error()
			}

			statusCode = jsonErr.Err.Code

			errb, _ := jsonErr.MarshalJSON()
			ctx.Error(string(errb), statusCode)
		}
	}
}
