package account

import (
	"context"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ravlio/wallet/account/pb"
)

type grpcGoKitServer struct {
	createAccount gt.Handler
	getAccount    gt.Handler
	deleteAccount gt.Handler
	updateAccount gt.Handler
	listAccounts  gt.Handler
}

func (g *grpcGoKitServer) CreateAccount(ctx context.Context, req *pb.Account) (resp *pb.Account, err error) {
	_, pbresp, err := g.createAccount.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return pbresp.(*pb.Account), nil
}

func (g *grpcGoKitServer) GetAccount(ctx context.Context, id *wrappers.UInt32Value) (*pb.Account, error) {
	_, pbresp, err := g.getAccount.ServeGRPC(ctx, id)
	if err != nil {
		return nil, err
	}

	return pbresp.(*pb.Account), nil
}

func (g *grpcGoKitServer) DeleteAccount(ctx context.Context, req *wrappers.UInt32Value) (*empty.Empty, error) {
	_, _, err := g.deleteAccount.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (g *grpcGoKitServer) UpdateAccount(ctx context.Context, req *pb.Account) (*pb.Account, error) {
	_, pbresp, err := g.updateAccount.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return pbresp.(*pb.Account), nil
}

func (g *grpcGoKitServer) ListAccount(_ *empty.Empty, ret pb.AccountService_ListAccountServer) error {
	// Я не знаю, как в go kit делать grpc streaming :/
	return nil
}

func (g *grpcGoKitServer) decodeCreateAccountRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.Account)
	return accountFromProtobuf(req), nil
}

func (g *grpcGoKitServer) encodeCreateAccountResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*Account)
	return resp.toProtobuf(), nil
}

func (g *grpcGoKitServer) decodeGetAccountRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*wrappers.UInt32Value)
	return req.Value, nil
}

func (g *grpcGoKitServer) encodeGetAccountResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*Account)
	return resp.toProtobuf(), nil
}

func (g *grpcGoKitServer) decodeDeleteAccountRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*wrappers.UInt32Value)
	return req.Value, nil
}

func (g *grpcGoKitServer) encodeDeleteAccountResponse(_ context.Context, r interface{}) (interface{}, error) {
	return &empty.Empty{}, nil
}

func (g *grpcGoKitServer) decodeUpdateAccountRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.Account)
	return accountFromProtobuf(req), nil
}

func (g *grpcGoKitServer) encodeUpdateAccountResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*Account)
	return resp.toProtobuf(), nil
}

func (g *grpcGoKitServer) decodeListAccountsRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, nil
}

func (g *grpcGoKitServer) encodeListAccountsResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.([]*Account)
	return req, nil
}

func NewGoKitGRPCServer(_ context.Context, ep Endpoints) pb.AccountServiceServer {
	srv := &grpcGoKitServer{}

	srv.createAccount = gt.NewServer(ep.CreateAccountEndpoint, srv.decodeCreateAccountRequest, srv.encodeCreateAccountResponse)
	srv.getAccount = gt.NewServer(ep.GetAccountEndpoint, srv.decodeGetAccountRequest, srv.encodeGetAccountResponse)
	srv.deleteAccount = gt.NewServer(ep.DeleteAccountEndpoint, srv.decodeDeleteAccountRequest, srv.encodeDeleteAccountResponse)
	srv.updateAccount = gt.NewServer(ep.UpdateAccountEndpoint, srv.decodeUpdateAccountRequest, srv.encodeUpdateAccountResponse)
	srv.listAccounts = gt.NewServer(ep.ListAccountsEndpoint, srv.decodeListAccountsRequest, srv.encodeListAccountsResponse)

	return srv
}
