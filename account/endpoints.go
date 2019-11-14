package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type CreateAccountRequest struct {
	Req *Account
}
type CreateAccountResponse struct {
	A   *Account
	Err error
}

func MakeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateAccountRequest)
		A, err := s.CreateAccount(ctx, req.Req)
		return CreateAccountResponse{A: A, Err: err}, nil
	}
}

type GetAccountRequest struct {
	Id uint32
}
type GetAccountResponse struct {
	A   *Account
	Err error
}

func MakeGetAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountRequest)
		A, err := s.GetAccount(ctx, req.Id)
		return GetAccountResponse{A: A, Err: err}, nil
	}
}

type DeleteAccountRequest struct {
	Id uint32
}
type DeleteAccountResponse struct {
	Err error
}

func MakeDeleteAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteAccountRequest)
		err := s.DeleteAccount(ctx, req.Id)
		return DeleteAccountResponse{Err: err}, nil
	}
}

type UpdateAccountRequest struct {
	Req *Account
}
type UpdateAccountResponse struct {
	A   *Account
	Err error
}

func MakeUpdateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateAccountRequest)
		A, err := s.UpdateAccount(ctx, req.Req)
		return UpdateAccountResponse{A: A, Err: err}, nil
	}
}

type ListAccountsRequest struct {
}
type ListAccountsResponse struct {
	S   []*Account
	Err error
}

func MakeListAccountsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListAccountsRequest)
		slice1, err := s.ListAccounts(ctx)
		return ListAccountsResponse{S: slice1, Err: err}, nil
	}
}

type Endpoints struct {
	CreateAccount endpoint.Endpoint
	GetAccount    endpoint.Endpoint
	DeleteAccount endpoint.Endpoint
	UpdateAccount endpoint.Endpoint
	ListAccounts  endpoint.Endpoint
}
