package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateAccountEndpoint endpoint.Endpoint
	GetAccountEndpoint    endpoint.Endpoint
	DeleteAccountEndpoint endpoint.Endpoint
	UpdateAccountEndpoint endpoint.Endpoint
	ListAccountsEndpoint  endpoint.Endpoint
}

func MakeEndpoints(svc Client) *Endpoints {
	return &Endpoints{
		CreateAccountEndpoint: func(ctx context.Context, req interface{}) (interface{}, error) {
			acc := req.(*Account)
			return svc.CreateAccount(ctx, acc)
		},
		GetAccountEndpoint: func(ctx context.Context, req interface{}) (interface{}, error) {
			id := req.(uint32)
			return svc.GetAccount(ctx, id)
		},
		DeleteAccountEndpoint: func(ctx context.Context, req interface{}) (interface{}, error) {
			id := req.(uint32)
			return nil, svc.DeleteAccount(ctx, id)
		},
		UpdateAccountEndpoint: func(ctx context.Context, req interface{}) (interface{}, error) {
			acc := req.(*Account)
			return svc.UpdateAccount(ctx, acc)
		},
		ListAccountsEndpoint: func(ctx context.Context, req interface{}) (interface{}, error) {
			return svc.ListAccounts(ctx)
		},
	}
}
