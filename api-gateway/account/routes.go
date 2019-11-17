package account

import (
	"github.com/ravlio/wallet/account"
	api_gateway "github.com/ravlio/wallet/api-gateway"
	"github.com/ravlio/wallet/api-gateway/middleware"
	"github.com/ravlio/wallet/pkg/errutil"
	"github.com/valyala/fasthttp"
)

type routes struct {
	cl account.Client
}

func New(r *api_gateway.Router, cl account.Client) {
	rr := &routes{cl: cl}

	r.POST("/accounts", middleware.Json(rr.createAccount))
}
func (r *routes) createAccount(ctx *fasthttp.RequestCtx) error {
	req := &Account{}
	if err := req.UnmarshalJSON(ctx.PostBody()); err != nil {
		return errutil.NewBadRequestError(err)
	}

	sresp, err := r.cl.CreateAccount(ctx, req.ToEntity())
	if err != nil {
		return err
	}

	resp := &Account{}
	resp.FromEntity(sresp)

	data, err := resp.MarshalJSON()
	if err != nil {
		return errutil.NewInternalServerError(err)
	}

	ctx.Response.SetStatusCode(fasthttp.StatusCreated)
	_, err = ctx.Write(data)

	return nil
}
