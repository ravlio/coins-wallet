package middleware

import "github.com/valyala/fasthttp"

type Endpoint func(ctx *fasthttp.RequestCtx) error
